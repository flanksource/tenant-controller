package utils

import (
	"bytes"
	"fmt"
	"math/rand"

	gotemplate "text/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func Template(template string, vars map[string]interface{}) (string, error) {
	tpl := gotemplate.New("")
	tpl, err := tpl.Parse(template)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RandomString(length uint) string {
	// Generate a random password of length 16
	const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	for i := range password {
		password[i] = passwordChars[rand.Intn(len(passwordChars))]
	}
	return string(password)
}

func GetUnstructuredObjects(data ...string) ([]*unstructured.Unstructured, error) {
	var items []*unstructured.Unstructured
	for _, d := range data {
		if d == "" {
			continue
		}
		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(d)), 1024)
		var resource *unstructured.Unstructured

		if err := decoder.Decode(&resource); err != nil {
			return nil, fmt.Errorf("error decoding %s: %s", d, err)
		}
		if resource != nil {
			items = append(items, resource)
		}
	}
	return items, nil
}
