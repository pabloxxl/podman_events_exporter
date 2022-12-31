package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pabloxxl/podman_events_exporter/pkg/events"
	"github.com/pabloxxl/podman_events_exporter/pkg/utils"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	klog "k8s.io/klog/v2"
)

var (
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

func loop(config *utils.ConfigOpts) {
	ctx, err := utils.ConnectToPodmanSocket(config.Socket)
	if err != nil {
		os.Exit(1)
	}

	run := true
	counters := make(map[string]*prometheus.CounterVec)

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
	go utils.CreateListener(ctx, &eventChan, &exitChan)

	klog.Infof("Listening on %s/metrics", config.HostWithPort)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(config.HostWithPort, nil)

	for run {
		msg := <-eventChan
		events.ConvertEventToCounter(&msg, counters, config.Include, config.Exclude, config.Regex)
	}
}
func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	config := utils.ParseCLIArguments()
	if config.Help {
		flag.Usage()
		os.Exit(0)
	} else if config.Version {
		printBuildData()
		os.Exit(0)
	}

	config.PrintParameters()

	// Disable ugly logs from podman library
	logrus.SetOutput(ioutil.Discard)

	klog.Infof("Running podman_events_exporter version %s", Version)
	loop(config)
	klog.V(2).Infof("Running podman_events_exporter version %s", Version)

	klog.Info("Program finished")

}
