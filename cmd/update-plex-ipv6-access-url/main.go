package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/update-plex-ipv6-access-url/cmd/update-plex-ipv6-access-url/internal/config"
	"github.com/cetteup/update-plex-ipv6-access-url/cmd/update-plex-ipv6-access-url/internal/handler"
	"github.com/cetteup/update-plex-ipv6-access-url/internal"
	"github.com/cetteup/update-plex-ipv6-access-url/internal/plex"
)

const (
	logKeyInterfaceName = "interfaceName"
	logKeyConfigPath    = "configPath"
)

var (
	buildVersion = "development"
	buildCommit  = "uncommitted"
	buildTime    = "unknown"
)

func main() {
	version := fmt.Sprintf("update-plex-ipv6-access-url %s (%s) built at %s", buildVersion, buildCommit, buildTime)
	cfg := config.Init()

	// Print version and exit
	if cfg.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:          os.Stdout,
		NoColor:      !cfg.ColorizeLogs,
		PartsExclude: []string{"time", "level"},
	})

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// TODO Validate cfg instead
	if cfg.InterfaceName == "" {
		log.Fatal().Msg("No interface name specified")
	}
	if cfg.ConfigPath == "" {
		log.Fatal().Msg("No config path specified")
	}

	a, err := internal.GetInterfaceGlobalUnicastIPv6ByName(cfg.InterfaceName)
	if err != nil {
		log.Fatal().
			Err(err).
			Str(logKeyInterfaceName, cfg.InterfaceName).
			Msg("Failed to find global unicast IPv6 address")
	}

	log.Info().
		Str(logKeyInterfaceName, cfg.InterfaceName).
		Str("address", a.String()).
		Msg("Found IPv6 address on interface")

	c, err := plex.ReadConfigFile(cfg.ConfigPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Str(logKeyConfigPath, cfg.ConfigPath).
			Msg("Failed to read Plex config file")
	}

	err = handler.UpdateIPv6CustomAccessURL(a, c)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to update custom access url")
	}

	log.Info().Msg("Successfully updated IPv6 custom server access URL")
}
