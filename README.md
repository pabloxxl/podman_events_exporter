# podman_events_exporter
Prometheus exporter that exports podman events


Example output from counters:
```
[root@host user]# curl -s localhost:2112/metrics | grep podman
# HELP podman_events_cleanup Podman event cleanup
# TYPE podman_events_cleanup counter
podman_events_cleanup{name="exciting_murdock"} 6
# HELP podman_events_create Podman event create
# TYPE podman_events_create counter
podman_events_create{name="gallant_hofstadter"} 1
podman_events_create{name="great_joliot"} 1
# HELP podman_events_died Podman event died
# TYPE podman_events_died counter
podman_events_died{name="exciting_murdock"} 3
# HELP podman_events_history Podman event history
# TYPE podman_events_history counter
podman_events_history{name="unkown"} 2
# HELP podman_events_init Podman event init
# TYPE podman_events_init counter
podman_events_init{name="exciting_murdock"} 3
podman_events_init{name="gallant_hofstadter"} 1
podman_events_init{name="great_joliot"} 1
# HELP podman_events_pull Podman event pull
# TYPE podman_events_pull counter
podman_events_pull{name="ubuntu"} 2
# HELP podman_events_start Podman event start
# TYPE podman_events_start counter
podman_events_start{name="exciting_murdock"} 3
podman_events_start{name="gallant_hofstadter"} 1
podman_events_start{name="great_joliot"} 1
# HELP podman_events_stop Podman event stop
# TYPE podman_events_stop counter
podman_events_stop{name="exciting_murdock"} 3
```
