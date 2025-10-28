package embed

import (
	"embed"
)

//go:embed locales/*.json
var FsLocales embed.FS
