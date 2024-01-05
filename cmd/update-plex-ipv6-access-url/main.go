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

	if err := cfg.ReadValuesIfMissing(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read missing config values")
	}

	interfaceAddrs, err := internal.GetGlobalUnicastIPv6AddrsByInterfaceName(cfg.InterfaceName)
	if err != nil {
		log.Fatal().
			Err(err).
			Str(logKeyInterfaceName, cfg.InterfaceName).
			Msg("Failed to find global unicast IPv6 addresses on interface")
	}

	if len(interfaceAddrs) == 0 {
		log.Fatal().
			Str(logKeyInterfaceName, cfg.InterfaceName).
			Msg("No global unicast IPv6 address found on interface")
	}

	log.Info().
		Str(logKeyInterfaceName, cfg.InterfaceName).
		Interface("addresses", interfaceAddrs).
		Msg("Found IPv6 addresses on interface")

	localClient := plex.NewApiClient(cfg.ServerAddr, cfg.Token, cfg.Timeout)
	remoteClient := plex.NewApiClient(plex.BaseURL, cfg.Token, cfg.Timeout)
	h := handler.NewHandler(localClient, remoteClient)

	selectedAddrs, err := h.SelectAddrs(interfaceAddrs, handler.AddrPreference(cfg.AddrPreference))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to select IPv6 addresses to use")
	}

	if len(interfaceAddrs) > 1 {
		log.Info().
			Stringer("use", cfg.AddrPreference).
			Interface("addresses", selectedAddrs).
			Msg("Selected IPv6 addresses")
	}

	err = h.UpdateIPv6CustomAccessURLs(selectedAddrs, handler.IPv6URLCapitalization(cfg.Capitalization))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to update custom access urls")
	}

	log.Info().Msg("Successfully updated IPv6 custom server access URLs")
}
