package domain

import "encoding/json"

// PatchInfo contains all information about a patch file
type PatchInfo struct {
	PatchPath  string            `json:"patchPath"`
	Dictionary map[string]string `json:"dictionary"`
	Overrides  []string          `json:"overrides"`
	Config     *Config           `json:"config"`
}

// Config defines patch configuration and rules
type Config struct {
	VariablesToPatch  []int              `json:"variablesToPatch"`
	WrapWidth         int                `json:"wrapWidth"`
	Version           int                `json:"version"`
	ParametersToPatch []ParameterToPatch `json:"parametersToPatch"`
	PluginsToPatch    []PluginToPatch    `json:"pluginsToPatch"`
	CreditsLocation   string             `json:"creditsLocation"`
	DynamicWrapWidth  bool               `json:"dynamicWrapWidth"`
	Locale            string             `json:"locale"`
}

// PluginToPatch defines how to patch a specific plugin
type PluginToPatch struct {
	Plugin                string              `json:"plugin"`
	ParametersPatchScript string              `json:"parametersPatchScript"` // Lua script
	ReplaceRules          []PluginReplaceRule `json:"replaceRules"`
}

// PluginReplaceRule defines a text replacement rule for plugin files
type PluginReplaceRule struct {
	Match   string `json:"match"`
	Replace string `json:"replace"`
}

// ParameterToPatch defines how to patch plugin command parameters
type ParameterToPatch struct {
	Plugin                string                 `json:"plugin"`
	Function              string                 `json:"function"`
	RootType              string                 `json:"rootType"`
	ParameterPathsToPatch []ParameterPathToPatch `json:"parameterPathsToPatch"`
}

// ParameterPathToPatch defines a specific parameter path to translate
type ParameterPathToPatch struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// PluginData represents a plugin entry in plugins.js
type PluginData struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      bool            `json:"status"`
	Parameters  json.RawMessage `json:"parameters"`
}
