package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
)

func OpenPatch(path string) (*zip.ReadCloser, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func readJSONFromZip[T any](r *zip.ReadCloser, name string) (*T, error) {
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			b, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			var v T
			return &v, json.Unmarshal(b, &v)
		}
	}
	return nil, errors.New(name + " not found")
}

func ReadDictionary(r *zip.ReadCloser) (map[string]string, error) {
	dictionary, err := readJSONFromZip[map[string]string](r, "dictionary.json")
	if err != nil {
		return nil, err
	}
	return *dictionary, nil
}

func ReadConfig(r *zip.ReadCloser) (*Config, error) {
	return readJSONFromZip[Config](r, "config.json")
}

type Config struct {
	VariablesToPatch  []int              `json:"variablesToPatch"`
	WrapWidth         int                `json:"wrapWidth"`
	Version           int                `json:"version"`
	ParametersToPatch []ParameterToPatch `json:"parametersToPatch"`
	PluginsToPatch    []PluginToPatch    `json:"pluginsToPatch"`
	CreditsLocation   string             `json:"creditsLocation"`
}

type PluginToPatch struct {
	Plugin                string              `json:"plugin"`
	ParametersPatchScript string              `json:"parametersPatchScript"` // this is a lua script
	ReplaceRules          []PluginReplaceRule `json:"replaceRules"`
}

type PluginReplaceRule struct {
	Match   string `json:"match"`
	Replace string `json:"replace"`
}

type ParameterToPatch struct {
	Plugin                string                 `json:"plugin"`
	Function              string                 `json:"function"`
	RootType              string                 `json:"rootType"`
	ParameterPathsToPatch []ParameterPathToPatch `json:"parameterPathsToPatch"`
}

type ParameterPathToPatch struct {
	Path string `json:"path"`
	Type string `json:"type"`
}
