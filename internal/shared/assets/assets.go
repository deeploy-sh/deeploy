package assets

import "embed"

//go:embed fonts/* js/* css/* img/*
var Assets embed.FS
