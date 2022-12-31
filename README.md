# podman_events_exporter
Prometheus exporter that exports podman events


Example output from counters:
```
[root@host user]# curl -s localhost:2112/metrics | grep podman
# HELP podman_events_cleanup Podman event cleanup
# TYPE podman_events_cleanup counter
podman_events_cleanup{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="awesome_mcclintock",network_backend="cni"} 1
# HELP podman_events_create Podman event create
# TYPE podman_events_create counter
podman_events_create{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="awesome_mcclintock",network_backend="cni"} 1
podman_events_create{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="determined_hermann",network_backend="cni"} 1
# HELP podman_events_died Podman event died
# TYPE podman_events_died counter
podman_events_died{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="awesome_mcclintock",network_backend="cni"} 1
# HELP podman_events_history Podman event history
# TYPE podman_events_history counter
podman_events_history{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",network_backend="cni"} 2
# HELP podman_events_init Podman event init
# TYPE podman_events_init counter
podman_events_init{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="awesome_mcclintock",network_backend="cni"} 1
podman_events_init{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="determined_hermann",network_backend="cni"} 1
# HELP podman_events_kill Podman event kill
# TYPE podman_events_kill counter
podman_events_kill{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="awesome_mcclintock",network_backend="cni"} 1
# HELP podman_events_pull Podman event pull
# TYPE podman_events_pull counter
podman_events_pull{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="ubuntu",network_backend="cni"} 2
# HELP podman_events_start Podman event start
# TYPE podman_events_start counter
podman_events_start{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="awesome_mcclintock",network_backend="cni"} 1
podman_events_start{api_version="4.1.1",arch="amd64",cgroups="systemd",go_version="go1.18.3",hostname="myhostname",name="determined_hermann",network_backend="cni"} 1
```

Example config.toml provided by `--config` option
```
Socket="/var/run/podman/podman.sock"
Host="127.0.0.1"
Port="2345"
Include=["kill","start"]
Exclude=["stop"]
Regex=".*a.*"
```

