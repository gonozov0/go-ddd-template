package ydb

import (
	"embed"
)

//go:embed *.sql
var FS embed.FS
