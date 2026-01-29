// Package schemas provides embedded JSON schema files
package schemas

import (
	"embed"
)

//go:embed *.json
var FS embed.FS
