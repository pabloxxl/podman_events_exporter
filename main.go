package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/system"
	"github.com/containers/podman/v4/pkg/domain/entities"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func create(ctx context.Context, eventChan chan entities.Event, exitChan chan bool) error {
	fmt.Println("Creating events listener...")
	err := system.Events(ctx, eventChan, exitChan, &system.EventsOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Events listener is finished...")
	return nil
}

func main() {

	sock_dir := "/var/run"
	socket := "unix:" + sock_dir + "/podman/podman.sock"

	fmt.Println("Connecting to podman socket...")
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(socket)

	exitChan := make(chan bool)
	eventChan := make(chan entities.Event)
	go create(ctx, eventChan, exitChan)

	exampleMetric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "podman_example_metric",
		Help: "Example podman metric",
	})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	exampleMetric.Inc()
	msg := <-eventChan
	fmt.Println(msg.ID)
	fmt.Println(msg.Action)
	fmt.Println(msg.Actor)
	exitChan <- true

	return

}
