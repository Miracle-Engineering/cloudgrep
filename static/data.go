package static

import "embed"

//go:embed img/* js/* css/*
//go:embed index.html
var Static embed.FS
