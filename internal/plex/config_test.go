package plex

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreferences_Validate(t *testing.T) {
	tests := []struct {
		name               string
		preparePreferences func(data Preferences)
		wantErrorContains  string
	}{
		{
			name:               "passes for valid preferences with manually mapped port",
			preparePreferences: func(data Preferences) {},
		},
		{
			name:               "passes for valid preferences with automatically mapped port",
			preparePreferences: func(data Preferences) {},
		},
		{
			name: "fails for missing MachineIdentifier",
			preparePreferences: func(data Preferences) {
				delete(data, "MachineIdentifier")
			},
			wantErrorContains: "failed on the 'required' tag",
		},
		{
			name: "fails for non-uuid4 MachineIdentifier",
			preparePreferences: func(data Preferences) {
				data["MachineIdentifier"] = "not-a-uuid4"
			},
			wantErrorContains: "failed on the 'uuid4' tag",
		},
		{
			name: "fails for missing ProcessedMachineIdentifier",
			preparePreferences: func(data Preferences) {
				delete(data, "ProcessedMachineIdentifier")
			},
			wantErrorContains: "failed on the 'required' tag",
		},
		{
			name: "fails for non-40 character ProcessedMachineIdentifier",
			preparePreferences: func(data Preferences) {
				data["ProcessedMachineIdentifier"] = strings.Repeat("a", 39)
			},
			wantErrorContains: "failed on the 'len' tag",
		},
		{
			name: "fails for non-hexadecimal ProcessedMachineIdentifier",
			preparePreferences: func(data Preferences) {
				data["ProcessedMachineIdentifier"] = strings.Repeat("z", 40)
			},
			wantErrorContains: "failed on the 'hexadecimal' tag",
		},
		{
			name: "fails for missing PlexOnlineToken",
			preparePreferences: func(data Preferences) {
				delete(data, "PlexOnlineToken")
			},
			wantErrorContains: "failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			configData := Preferences{
				preferenceKeyMachineIdentifier:          "e19d8db6-bb1c-45fa-8cb4-63116df5c8e1",
				preferenceKeyProcessedMachineIdentifier: "1142ed040a27acc36ea876e8362b28464c3d240d",
				preferenceKeyPlexOnlineToken:            "Ydu_4wQEX7Lmt4HPDy8A",
				preferenceKeyManualPortMappingMode:      "1",
				preferenceKeyManualPortMappingPort:      "0",
			}
			tt.preparePreferences(configData)

			// WHEN
			err := configData.Validate()

			// THEN
			if tt.wantErrorContains != "" {
				assert.ErrorContains(t, err, tt.wantErrorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name              string
		givenData         []byte
		wantPreferences   Preferences
		wantErrorContains string
	}{
		{
			name:      "successfully reads plex config file",
			givenData: []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<Preferences MachineIdentifier=\"e19d8db6-bb1c-45fa-8cb4-63116df5c8e1\" FriendlyName=\"MyPlexServer\" customConnections=\"\" EnableIPv6=\"1\"/>"),
			wantPreferences: Preferences{
				"MachineIdentifier": "e19d8db6-bb1c-45fa-8cb4-63116df5c8e1",
				"FriendlyName":      "MyPlexServer",
				"customConnections": "",
				"EnableIPv6":        "1",
			},
		},
		{
			name:      "does not error if prolog is missing",
			givenData: []byte("<Preferences FriendlyName=\"MyPlexServer\"/>"),
			wantPreferences: Preferences{
				"FriendlyName": "MyPlexServer",
			},
		},
		{
			name:            "does not error if preferences tag does not have any attributes",
			givenData:       []byte("<Preferences/>"),
			wantPreferences: Preferences{},
		},
		{
			name:              "returns error if preferences tag is not opened",
			givenData:         []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\nPreferences FriendlyName=\"MyPlexServer\"/>"),
			wantErrorContains: "EOF",
		},
		{
			name:              "returns error if preferences tag is not closed",
			givenData:         []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<Preferences FriendlyName=\"MyPlexServer\""),
			wantErrorContains: "XML syntax error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			f, err := os.CreateTemp(os.TempDir(), "Preferences.xml")
			require.NoError(t, err)

			t.Cleanup(func() {
				_ = os.Remove(f.Name())
			})

			_, err = f.Write(tt.givenData)
			require.NoError(t, err)
			err = f.Close()
			require.NoError(t, err)

			// WHEN
			config, err := ReadConfigFile(f.Name())

			// THEN
			if tt.wantErrorContains != "" {
				assert.ErrorContains(t, err, tt.wantErrorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, config.Path, f.Name())
				assert.Equal(t, tt.wantPreferences, config.Preferences)
			}
		})
	}
}

func TestWriteConfigFile(t *testing.T) {
	tests := []struct {
		name             string
		givenPreferences Preferences
		wantDataElements []string
	}{
		{
			name: "successfully writes plex config file",
			givenPreferences: Preferences{
				"MachineIdentifier": "e19d8db6-bb1c-45fa-8cb4-63116df5c8e1",
				"FriendlyName":      "MyPlexServer",
				"customConnections": "",
				"EnableIPv6":        "1",
			},
			wantDataElements: []string{
				"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<Preferences",
				"MachineIdentifier=\"e19d8db6-bb1c-45fa-8cb4-63116df5c8e1\"",
				"FriendlyName=\"MyPlexServer\"",
				"customConnections=\"\"",
				"EnableIPv6=\"1\"",
				"></Preferences>",
			},
		},
		{
			name:             "does not error for empty config",
			givenPreferences: Preferences{},
			wantDataElements: []string{"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<Preferences></Preferences>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			f, err := os.CreateTemp(os.TempDir(), "Preferences.xml")
			require.NoError(t, err)

			t.Cleanup(func() {
				_ = os.Remove(f.Name())
			})

			err = f.Close()
			require.NoError(t, err)

			config := Config{
				Path:        f.Name(),
				Preferences: tt.givenPreferences,
			}

			// WHEN
			err = WriteConfigFile(config)

			// THEN
			assert.NoError(t, err)
			data, err := os.ReadFile(f.Name())
			require.NoError(t, err)
			for _, elem := range tt.wantDataElements {
				assert.Contains(t, string(data), elem)
			}
		})
	}
}
