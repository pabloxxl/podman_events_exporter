package events

import (
	"regexp"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog"
)

func ConvertEventToCounter(event *entities.Event, counters map[string]*prometheus.CounterVec,
	include map[string]bool, exclude map[string]bool, rgx *regexp.Regexp) {

	val, ok := event.Actor.Attributes["name"]
	name := "unkown"
	action := event.Action
	var labelNames []string
	labels := make(map[string]string)

	if ok && val != "" {
		name = val
		labels["name"] = name
		labelNames = append(labelNames, "name")
	}

	if rgx != nil {
		if !rgx.Match([]byte(name)) {
			klog.V(3).Infof("Dropping %s for %s: regular expression does not match", action, name)
			return
		}

	}
	if len(include) > 0 && !include[action] {
		klog.V(3).Infof("Dropping %s for %s: action is not included", action, name)
		return
	}

	if len(exclude) > 0 && exclude[action] {
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
