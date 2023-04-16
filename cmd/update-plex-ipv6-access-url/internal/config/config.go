package config

import (
	"flag"
)

type Config struct {
	Version bool

	Debug        bool
	ColorizeLogs bool

	InterfaceName string
	ConfigPath    string
}

func Init() *Config {
	cfg := new(Config)
	flag.BoolVar(&cfg.Version, "v", false, "prints the version")
	flag.BoolVar(&cfg.Version, "version", false, "prints the version")
	flag.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
	flag.BoolVar(&cfg.ColorizeLogs, "colorize-logs", false, "colorize log messages")
	flag.StringVar(&cfg.InterfaceName, "interface", "", "Name of network interface to use for IPv6 access")
	flag.StringVar(&cfg.ConfigPath, "config", "", "Path to Plex config (Preferences.xml)")
	flag.Parse()
	return cfg
}
