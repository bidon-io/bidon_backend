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
