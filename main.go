package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/system"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	klog "k8s.io/klog/v2"
)

const (
	SOCK_DIR = "/var/run"
)

var (
	argSocket          = flag.String("socket", "", "Podman socket path")
	Version     string = "n/a"
	BuildCommit string = "n/a"
	BuildBranch string = "n/a"
	BuildHost   string = "n/a"
	BuildTime   string = "n/a"
)

func printBuildData() {
	fmt.Printf("Build variables of podman_events_exporter_%s:\n", Version)
	fmt.Printf("  commit:     %s\n", BuildCommit)
	fmt.Printf("  branch:     %s\n", BuildBranch)
	fmt.Printf("  build host: %s\n", BuildHost)
	fmt.Printf("  build time: %s\n", BuildTime)
}

func parseArguments() (*bool, *bool, *string, map[string]bool, map[string]bool) {
	include := make(map[string]bool)
	exclude := make(map[string]bool)

	argVersion := flag.Bool("version", false, "Print version and exit")
	arghelp := flag.Bool("help", false, "Print help and exit")
	argHost := flag.String("host", "", "Host to serve metrics on")
	argPort := flag.String("port", "2112", "Port to serve metrics on")
	argInclude := flag.String("include", "", "Include certain events, comma separated")
	argExclude := flag.String("exclude", "", "Exclude certain events, comma separated")

	flag.Parse()
	for _, elem := range strings.Split(*argInclude, ",") {
		if len(elem) > 2 {
			include[elem] = true
		}
	}

	for _, elem := range strings.Split(*argExclude, ",") {
		if len(elem) > 2 {
			exclude[elem] = true
		}
	}
	hostWithPort := *argHost + ":" + *argPort

	return argVersion, arghelp, &hostWithPort, include, exclude
}

func createListener(ctx context.Context, eventChan *chan entities.Event, exitChan *chan bool) error {
	klog.Info("Creating events listener")
	err := system.Events(ctx, *eventChan, *exitChan, &system.EventsOptions{})
	if err != nil {
		klog.V(2).ErrorS(err, "Event is missing action type")
	}
	klog.Info("Events listener is finished")
	return nil
}

func convertEventToCounter(event *entities.Event, counters map[string]*prometheus.CounterVec, include map[string]bool, exclude map[string]bool) {
	val, ok := event.Actor.Attributes["name"]
	name := "unkown"
	action := event.Action
	var labelNames []string
	labels := make(map[string]string)

	if len(include) > 0 && !include[action] {
		klog.V(2).Infof("%s is not included. Included labels: %s", action, include)
		return
	}

	if len(exclude) > 0 && include[action] {
		klog.V(2).Infof("%s is excluded", action)
		return
	}

	if action == "" {
		klog.V(2).Info("Event is missing action type")
		return
	}

	if ok && val != "" {
		name = val
		labels["name"] = name
		labelNames = append(labelNames, "name")
	}

	valC, okC := counters[action]
	if !okC {
		klog.V(2).Infof("Creating new counter: podman_events_%s with %d labels", event.Action, len(labelNames))
		counters[event.Action] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "podman_events_" + action,
				Help: "Podman event " + action,
			},
			labelNames)
		valC = counters[action]
		prometheus.MustRegister(valC)
	}

	klog.V(2).Infof("Incrementing counter: podman_events_%s for %s with labels %s", action, name, labels)
	valC.With(labels).Inc()
}

func connectToPodmanSocket(path string) (context.Context, error) {
	socket := "unix:" + path + "/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		klog.Errorf("Failed to connect to %s", socket)
		return nil, err
	}
	klog.Infof("Connected to podman socket at %s", socket)
	return ctx, nil
}

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	version, help, hostWithPort, include, exclude := parseArguments()
	if *help {
		flag.Usage()
		os.Exit(0)
	} else if *version {
		printBuildData()
		os.Exit(0)
	}
	// Disable ugly logs from podman library
	logrus.SetOutput(ioutil.Discard)

	klog.Infof("Running podman_events_exporter version %s", Version)
	klog.V(2).Infof("Running podman_events_exporter version %s", Version)
	for k := range include {
		klog.Infof("Including event %s", k)
	}

	for k := range exclude {
		klog.Infof("Excluding event %s", k)
	}

	run := true
	counters := make(map[string]*prometheus.CounterVec)

	sock_dir := SOCK_DIR
	if *argSocket != "" {
		sock_dir = *argSocket
	}

	ctx, err := connectToPodmanSocket(sock_dir)
	if err != nil {
		os.Exit(1)
	}

	exitChan := make(chan bool)
	eventChan := make(chan entities.Event)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		klog.Warningf("Caught signal: %s", sig.String())
		exitChan <- true
		run = false
	}()
	go createListener(ctx, &eventChan, &exitChan)

	klog.Infof("Listening on %s/metrics", *hostWithPort)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(*hostWithPort, nil)

	for run {
		msg := <-eventChan
		convertEventToCounter(&msg, counters, include, exclude)
	}

	klog.Info("Program finished")

}
