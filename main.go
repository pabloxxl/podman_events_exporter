package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/system"
	"github.com/containers/podman/v4/pkg/domain/entities"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	klog "k8s.io/klog/v2"
)

const (
	SOCK_DIR = "/var/run"
)

func createListener(ctx context.Context, eventChan *chan entities.Event, exitChan *chan bool) error {
	klog.Info("Creating events listener")
	err := system.Events(ctx, *eventChan, *exitChan, &system.EventsOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	klog.Info("Events listener is finished")
	return nil
}

func convertEventToCounter(event *entities.Event, counters map[string]*prometheus.CounterVec) {
	val, ok := event.Actor.Attributes["name"]
	name := "unkown"
	action := event.Action

	if ok && val != "" {
		name = val
	}

	valC, okC := counters[action]
	if !okC {
		klog.Infof("Creating new counter: podman_events_%s", event.Action)
		counters[event.Action] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "podman_events_" + action,
				Help: "Podman event " + action,
			},
			[]string{"name"})
		valC = counters[action]
		prometheus.MustRegister(valC)
	}

	klog.Infof("Incrementing counter: podman_events_%s for %s", action, name)
	valC.With(prometheus.Labels{"name": name}).Inc()
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

	flag.Set("alsologtostderr", "true")
	flag.Parse()
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	run := true
	counters := make(map[string]*prometheus.CounterVec)

	ctx, err := connectToPodmanSocket(SOCK_DIR)
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

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	for run {
		msg := <-eventChan
		convertEventToCounter(&msg, counters)
	}

	klog.Info("Program finished")
	return

}
