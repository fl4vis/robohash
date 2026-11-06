// Package assets provides the images compiled into the application binary.
package assets

import "embed"

// Files contains every generator asset.
//
//go:embed img
var Files embed.FS
