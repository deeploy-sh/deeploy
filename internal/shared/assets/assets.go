package assets

import "embed"

//go:embed fonts/* js/* css/* img/* video/*
var Assets embed.FS
