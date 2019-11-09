package main

import (
	ctx "context"

	"github.com/go-pg/pg/v9"
)

type dbLogger struct {
	log zerolog.Logger
}

func newDBLogger(log zerolog.Logger) dbLogger {
	return dbLogger{log: log}
}

func (d dbLogger) BeforeQuery(c ctx.Context, q *pg.QueryEvent) (ctx.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c ctx.Context, q *pg.QueryEvent) error {
	qs, err := q.FormattedQuery()
	if err != nil {
		d.log.Error().Err(err).Msg("query string error)
		return err
	}
	d.log.Debug().Str("query", qs).Msg("database query debug")
	return nil
}
