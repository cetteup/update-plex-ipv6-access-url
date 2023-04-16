package plex

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	xmlHeader = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"

	preferenceKeyMachineIdentifier          = "MachineIdentifier"
	preferenceKeyProcessedMachineIdentifier = "ProcessedMachineIdentifier"
	preferenceKeyPlexOnlineToken            = "PlexOnlineToken"
	preferenceKeyManualPortMappingMode      = "ManualPortMappingMode"
	preferenceKeyManualPortMappingPort      = "ManualPortMappingPort"
	preferenceKeyLastAutomaticMappedPort    = "LastAutomaticMappedPort"
	preferenceKeyCustomConnections          = "customConnections"
)

type Config struct {
	Path        string
	Preferences Preferences
}

type Preferences map[string]string

func (p *Preferences) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	*p = Preferences{}

	for _, attr := range start.Attr {
		(*p)[attr.Name.Local] = attr.Value
	}

	_, err := decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func (p *Preferences) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	attrs := make([]xml.Attr, 0, len(*p))
	for k, v := range *p {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: k},
			Value: v,
		})
	}
	start.Attr = attrs

	// Use EndElement, since go does not (yet) support self-closing tags, see https://github.com/golang/go/issues/21399
	tokens := []xml.Token{start, start.End()}
	for _, t := range tokens {
		err := encoder.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	return encoder.Flush()
}

func (p *Preferences) Validate() error {
	validate := validator.New()

	data := make(map[string]interface{})
	for k, v := range *p {
		data[k] = v
	}
	rules := map[string]interface{}{
		preferenceKeyMachineIdentifier:          "required,uuid4",
		preferenceKeyProcessedMachineIdentifier: "required,len=40,hexadecimal",
		preferenceKeyPlexOnlineToken:            "required",
		preferenceKeyManualPortMappingMode:      "required,oneof=0 1",
		preferenceKeyManualPortMappingPort:      "omitempty,numeric",
		preferenceKeyLastAutomaticMappedPort:    "omitempty,numeric",
	}

	if errs := validate.ValidateMap(data, rules); len(errs) > 0 {
		for _, e := range errs {
			validationErrs, ok := e.(validator.ValidationErrors)
			if ok {
				return validationErrs
			}
		}
		return fmt.Errorf("plex config data validation failed")
	}

	return nil
}

func (p *Preferences) getValue(key string) string {
	return (*p)[key]
}

func (p *Preferences) GetToken() string {
	return p.getValue(preferenceKeyPlexOnlineToken)
}

func (p *Preferences) GetProcessedMachineIdentifier() string {
	return p.getValue(preferenceKeyProcessedMachineIdentifier)
}

func (p *Preferences) GetMappedPort() string {
	if p.getValue(preferenceKeyManualPortMappingMode) == plexTrue {
		return p.getValue(preferenceKeyManualPortMappingPort)
	}
	return p.getValue(preferenceKeyLastAutomaticMappedPort)
}

func (p *Preferences) GetCustomConnections() []string {
	return strings.Split(p.getValue(preferenceKeyCustomConnections), ",")
}

func ReadConfigFile(path string) (Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var data Preferences
	if err := xml.Unmarshal(bytes, &data); err != nil {
		return Config{}, err
	}

	return Config{
		Path:        path,
		Preferences: data,
	}, nil
}

func WriteConfigFile(config Config) error {
	bytes, err := xml.Marshal(&config.Preferences)
	if err != nil {
		return err
	}

	// Marshal-ing does not add a prolog/header, so add it now (followed by a line-break)
	bytes = append([]byte(fmt.Sprintf("%s\n", xmlHeader)), bytes...)

	return os.WriteFile(config.Path, bytes, 0600)
}
