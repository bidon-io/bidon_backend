package admin

type App struct {
	ID int64 `json:"id"`
	AppAttrs
}

type AppAttrs struct {
	PlatformID  PlatformID     `json:"platform_id"`
	HumanName   string         `json:"human_name"`
	PackageName string         `json:"package_name"`
	UserID      int64          `json:"user_id"`
	AppKey      string         `json:"app_key"`
	Settings    map[string]any `json:"settings"`
}

type PlatformID string

const (
	UnknownPlatformID PlatformID = ""
	IOSPlatformID     PlatformID = "ios"
	AndroidPlatformID PlatformID = "android"
)

type AppService = resourceService[App, AppAttrs]
