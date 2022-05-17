package gen

import (
	"embed"
	_ "embed"
)

//go:embed template/*
var defaultTemplateFS embed.FS
