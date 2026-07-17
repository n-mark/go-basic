package migrations

import "embed"

// FS содержит встроенные SQL-файлы миграций.
// Используется goose через goose.SetBaseFS(migrations.FS)
// для прогона миграций из бинарника без обращения к файловой системе.
//
//go:embed *.sql
var FS embed.FS
