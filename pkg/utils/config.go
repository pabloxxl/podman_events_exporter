package utils

import (
	"flag"
	"os"
	"regexp"
	"strings"

	klog "k8s.io/klog/v2"
)

const (
	SOCK_DIR = "/var/run"
)

type ConfigOpts struct {
	Version      bool
	Help         bool
	Socket       string
	HostWithPort string
	Include      map[string]bool
	Exclude      map[string]bool
	Regex        *regexp.Regexp
}

func ParseCLIArguments() *ConfigOpts {
	include := make(map[string]bool)
	exclude := make(map[string]bool)

	klog.InitFlags(nil)
	defer klog.Flush()

	argVersion := flag.Bool("version", false, "Print version and exit")
	arghelp := flag.Bool("help", false, "Print help and exit")
	argSocket := flag.String("socket", "", "Podman socket path")
	argHost := flag.String("host", "", "Host to serve metrics on")
	argPort := flag.String("port", "2112", "Port to serve metrics on")
	argInclude := flag.String("include", "", "Include certain events, comma separated")
	argExclude := flag.String("exclude", "", "Exclude certain events, comma separated")
	argContainerRegex := flag.String("container_regex", "", "Container regular expression")

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

	socket := SOCK_DIR
	if *argSocket != "" {
		socket = *argSocket
	}

	config := ConfigOpts{*argVersion, *arghelp, socket, hostWithPort, include, exclude, nil}
	config.checkRegex(*argContainerRegex)

	return &config
}

func (c *ConfigOpts) PrintParameters() {
	for k := range c.Include {
		klog.Infof("Including event %s", k)
	}

	for k := range c.Exclude {
		klog.Infof("Excluding event %s", k)
	}

}

func (c *ConfigOpts) checkRegex(rawRegex string) {
	var err error
	if rawRegex != "" {
		c.Regex, err = regexp.Compile(rawRegex)
		if err != nil {
			klog.Errorf("%s is not a valid regular expression", rawRegex)
			os.Exit(1)
		}
	}
}
