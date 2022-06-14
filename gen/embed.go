package gen

import (
	"embed"
)

//go:embed template/*
var defaultTemplateFS embed.FS
