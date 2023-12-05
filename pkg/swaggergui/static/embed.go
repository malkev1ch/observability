// Package static contains files to embed.
package static

import (
	"embed"
)

// FS holds embedded static assets.
//
//go:embed *.png *.gz
var FS embed.FS
