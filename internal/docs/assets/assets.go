package assets

import "embed"

//go:embed css/* fonts/* js/*
var Assets embed.FS
