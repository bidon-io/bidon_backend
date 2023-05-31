package web

import "embed"

//go:embed redoc admin_api.yml
var FS embed.FS
