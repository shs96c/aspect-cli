package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// AspectPlugin represents a plugin entry in the plugins file.
type AspectPlugin struct {
	Name       string                 `yaml:"name"`
	From       string                 `yaml:"from"`
	LogLevel   string                 `yaml:"log_level"`
	Properties map[string]interface{} `yaml:"properties"`
}

// Parser is the interface that wraps the Parse method that performs the parsing
// of a plugins file.
type Parser interface {
	Parse(aspectpluginsPath string) ([]AspectPlugin, error)
}

type parser struct {
	ioutilReadFile      func(filename string) ([]byte, error)
	yamlUnmarshalStrict func(in []byte, out interface{}) (err error)
	yamlMarshal         func(in interface{}) (out []byte, err error)
}

// NewParser instantiates a default internal implementation of the Parser
// interface.
func NewParser() Parser {
	return &parser{
		ioutilReadFile:      ioutil.ReadFile,
		yamlUnmarshalStrict: yaml.UnmarshalStrict,
		yamlMarshal:         yaml.Marshal,
	}
}

// Parse parses a plugins file.
func (p *parser) Parse(aspectpluginsPath string) ([]AspectPlugin, error) {
	if aspectpluginsPath == "" {
		return []AspectPlugin{}, nil
	}
	aspectpluginsData, err := p.ioutilReadFile(aspectpluginsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse .aspectplugins: %w", err)
	}
	var aspectplugins []AspectPlugin
	if err := p.yamlUnmarshalStrict(aspectpluginsData, &aspectplugins); err != nil {
		return nil, fmt.Errorf("failed to parse .aspectplugins: %w", err)
	}

	return aspectplugins, nil
}
