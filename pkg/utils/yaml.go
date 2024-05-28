package utils

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

func MarshalYAML(obj any) ([]byte, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	// Use 2 space indentation to be consistent with yq
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(obj)
	return b.Bytes(), err
}
