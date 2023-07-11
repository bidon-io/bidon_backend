package schema

type App struct {
	Bundle           string   `json:"bundle" validate:"required"`
	Key              string   `json:"key" validate:"required"`
	Framework        string   `json:"framework" validate:"required"`
	Version          string   `json:"version" validate:"required"`
	FrameworkVersion string   `json:"framework_version"`
	PluginVersion    string   `json:"plugin_version"`
	Skadn            []string `json:"skadn"`
}

func (a App) Map() map[string]any {
	m := map[string]any{
		"bundle":            a.Bundle,
		"key":               a.Key,
		"framework":         a.Framework,
		"version":           a.Version,
		"framework_version": a.FrameworkVersion,
		"plugin_version":    a.PluginVersion,
		"skadn":             a.Skadn,
	}

	return m
}
