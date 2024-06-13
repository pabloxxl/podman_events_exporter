package events

import (
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/pabloxxl/podman_events_exporter/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	klog "k8s.io/klog/v2"
)

func ConvertEventToCounter(event *entities.Event, counters map[string]*prometheus.CounterVec,
	config *utils.ConfigOpts, infoLabels map[string]string) {

	val, ok := event.Actor.Attributes["name"]
	name := "unkown"
	action := event.Action
	var labelNames []string
	labels := make(map[string]string)

	for k, v := range infoLabels {
		labels[k] = v
		labelNames = append(labelNames, k)
	}

	if ok && val != "" {
		name = val
		labels["name"] = name
		labelNames = append(labelNames, "name")
	} else {
		// If podman socket dies, it might start sending empty events. The best way to handle this is to panic.
		panic(fmt.Sprintf("Missing or empty 'name' attribute in event.Actor.Attributes: %+v", event.Actor.Attributes))
	}

	if config.Regex != nil {
		if !config.Regex.Match([]byte(name)) {
			klog.V(3).Infof("Dropping %s for %s: regular expression does not match", action, name)
			return
		}

	}

	if len(config.Include) > 0 && !config.Include[action] {
		klog.V(3).Infof("Dropping %s for %s: action is not included", action, name)
		return
	}

	if len(config.Exclude) > 0 && config.Exclude[action] {
		klog.V(3).Infof("Dropping %s for %s: action is excluded", action, name)
		return
	}

	if action == "" {
		klog.V(2).Infof("Dropping %s for %s: missing action type", action, name)
		return
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
