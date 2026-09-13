package main

import "embed"

// embeddedFrontend holds the built frontend (GAHFrontend/dist) so the server can
// serve the web UI straight from the binary.
//
// The embed pattern requires GAHFrontend/dist to exist at compile time, which is
// why GAHFrontend/dist/.gitkeep is tracked in git. Run `npm run build` inside
// GAHFrontend before `go build` to embed the real UI instead of the placeholder.
//
//go:embed all:GAHFrontend/dist
var embeddedFrontend embed.FS
