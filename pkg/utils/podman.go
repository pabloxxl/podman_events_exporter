package utils

import (
	"context"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/system"
	"github.com/containers/podman/v4/pkg/domain/entities"
	klog "k8s.io/klog/v2"
)

func CreateListener(ctx context.Context, eventChan *chan entities.Event, exitChan *chan bool, breakChan *chan bool) error {
	klog.Info("Creating events listener")
	err := system.Events(ctx, *eventChan, *exitChan, &system.EventsOptions{})
	if err != nil {
		klog.V(2).ErrorS(err, "Event is missing action type")
	}
	klog.Info("Events listener is finished")
	*breakChan <- true
	return nil
}

func ConnectToPodmanSocket(path string) (context.Context, error) {
	socket := "unix:" + path
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		klog.Errorf("Failed to connect to %s", socket)
		return nil, err
	}
	klog.Infof("Connected to podman socket at %s", socket)
	return ctx, nil
}

func GetInfoLabels(ctx context.Context) (map[string]string, error) {
	infoLabels := make(map[string]string)

	info, err := system.Info(ctx, nil)
	if err != nil {
		return infoLabels, err
	}

	infoLabels["api_version"] = info.Version.APIVersion
	infoLabels["go_version"] = info.Version.GoVersion
	infoLabels["arch"] = info.Host.Arch
	infoLabels["cgroups"] = info.Host.CgroupManager
	infoLabels["hostname"] = info.Host.Hostname
	infoLabels["network_backend"] = info.Host.NetworkBackend

	return infoLabels, nil
}
