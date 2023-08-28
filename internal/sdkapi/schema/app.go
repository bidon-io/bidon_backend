package schema

type App struct {
	Bundle           string   `json:"bundle" validate:"required"`
	Key              string   `json:"key" validate:"required"`
	Framework        string   `json:"framework" validate:"required"`
	Version          string   `json:"version" validate:"required"`
	FrameworkVersion string   `json:"framework_version"`
	PluginVersion    string   `json:"plugin_version"`
	SKAdN            []string `json:"skadn"`
	SDKVersion       string   `json:"sdk_version"`
}

func (a App) Map() map[string]any {
	m := map[string]any{
		"bundle":            a.Bundle,
		"key":               a.Key,
		"framework":         a.Framework,
		"version":           a.Version,
		"framework_version": a.FrameworkVersion,
		"plugin_version":    a.PluginVersion,
		"sdk_version":       a.SDKVersion,
	}

	if a.SKAdN != nil {
		m["skadn"] = a.SKAdN
	}

	return m
}
