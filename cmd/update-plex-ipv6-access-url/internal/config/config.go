package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/cetteup/update-plex-ipv6-access-url/cmd/update-plex-ipv6-access-url/internal/handler"
	"github.com/cetteup/update-plex-ipv6-access-url/internal/plex"
)

type Config struct {
	Version bool

	Debug        bool
	ColorizeLogs bool

	ServerAddr     string
	InterfaceName  string
	AddrPreference handler.AddrPreference
	ConfigPath     string
	Token          string
	Capitalization handler.IPv6URLCapitalization
	Timeout        int
}

func Init() *Config {
	cfg := new(Config)
	flag.BoolVar(&cfg.Version, "v", false, "prints the version")
	flag.BoolVar(&cfg.Version, "version", false, "prints the version")
	flag.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
	flag.BoolVar(&cfg.ColorizeLogs, "colorize-logs", false, "colorize log messages")
	flag.StringVar(&cfg.ServerAddr, "address", "", "Plex server's address in format http[s]://host:port")
	flag.StringVar(&cfg.InterfaceName, "interface", "", "Name of network interface to use for IPv6 access")
	flag.TextVar(&cfg.AddrPreference, "use", handler.AddrPreferenceFirst, "Which IPv6 address(es) to use if multiple are found on the interface (first|last|all)")
	flag.StringVar(&cfg.ConfigPath, "config", "", "Path to Plex config (Preferences.xml)")
	flag.StringVar(&cfg.Token, "token", "", "Plex access token (X-Plex-Token) [required if 'config' flag is/cannot be provided]")
	flag.TextVar(&cfg.Capitalization, "capitalization", handler.IPv6URLCapitalizationLower, "Capitalization to use for dashed IPv6 address in Plex custom access URL (upper|lower)")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "Plex API request timeout (in seconds)")
	flag.Parse()
	return cfg
}

func (c *Config) ReadValuesIfMissing() error {
	if c.ServerAddr == "" {
		serverAddr, err := getInput("Enter the Plex server's address in format 'http[s]://host:port'")
		if err != nil {
			return fmt.Errorf("failed to read server address from console: %w", err)
		}
		c.ServerAddr = serverAddr
	}

	if c.InterfaceName == "" {
		interfaceName, err := getInput("Enter the name of network interface to use for IPv6 access")
		if err != nil {
			return fmt.Errorf("failed to read interface name from console: %w", err)
		}
		c.InterfaceName = interfaceName
	}

	if c.ConfigPath == "" && c.Token == "" {
		token, err := getInput("Enter a Plex access token (X-Plex-Token)")
		if err != nil {
			return fmt.Errorf("failed to read Plex token from console: %w", err)
		}
		c.Token = token
	}

	if c.ConfigPath != "" {
		config, err := plex.ReadConfigFile(c.ConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read Plex config file from %s: %w", c.ConfigPath, err)
		}

		c.Token = config.Preferences.GetToken()
	}

	return nil
}

func getInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRightFunc(input, func(r rune) bool {
		return unicode.IsControl(r)
	}), nil
}
