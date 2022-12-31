package utils

import (
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	klog "k8s.io/klog/v2"
)

const (
	SOCK_DIR     = "/var/run/podman/podman.sock"
	DEFAULT_PORT = "2112"
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

type FileConfigOpts struct {
	Socket  string
	Host    string
	Port    string
	Include []string
	Exclude []string
	Regex   string
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
	argPort := flag.String("port", "", "Port to serve metrics on")
	argInclude := flag.String("include", "", "Include certain events, comma separated")
	argExclude := flag.String("exclude", "", "Exclude certain events, comma separated")
	argContainerRegex := flag.String("container_regex", "", "Container regular expression")

	argsConfig := flag.String("config", "", "Config file path")

	flag.Parse()

	if *argsConfig != "" {
		fileConfig, err := parseConfigFile(*argsConfig)
		if err != nil {
			klog.Errorf("Cannot open provided config file %s: %s", *argsConfig, err)
			os.Exit(1)
		}
		argSocket = &fileConfig.Socket
		argHost = &fileConfig.Host
		argPort = &fileConfig.Port
		for _, elem := range fileConfig.Include {
			include[elem] = true
		}

		for _, elem := range fileConfig.Exclude {
			exclude[elem] = true
		}

		argContainerRegex = &fileConfig.Regex
	} else {

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
	}

	socket := SOCK_DIR
	if *argSocket != "" {
		socket = *argSocket
	}

	port := DEFAULT_PORT
	if *argPort != "" {
		port = *argPort
	}

	hostWithPort := *argHost + ":" + port

	config := ConfigOpts{*argVersion, *arghelp, socket, hostWithPort, include, exclude, nil}
	config.checkRegex(*argContainerRegex)

	return &config
}

func parseConfigFile(path string) (*FileConfigOpts, error) {
	var config FileConfigOpts
	_, err := toml.DecodeFile(path, &config)
	return &config, err
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
