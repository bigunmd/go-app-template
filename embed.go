package goapptemplate

import "embed"

//go:embed migrations/app
var MigrationsApp embed.FS
