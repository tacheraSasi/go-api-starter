package static

import (
	"embed"
)

//go:embed all:frontend/dist
var FrontendFiles embed.FS

func init() {
	//TODO: im gonna do something here later not sure what yet
}
