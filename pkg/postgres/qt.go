package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type queryTracer struct {
	log *logrus.Entry
}

// TraceQueryEnd implements pgx.QueryTracer.
func (t *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	t.log.WithFields(
		map[string]interface{}{
			"host":     conn.Config().Host,
			"port":     conn.Config().Port,
			"user":     conn.Config().User,
			"database": conn.Config().Database,
			"err":      data.Err,
		},
	).Debug("query execution end")
}

// TraceQueryStart implements pgx.QueryTracer.
func (t *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.log.WithFields(
		map[string]interface{}{
			"host":     conn.Config().Host,
			"port":     conn.Config().Port,
			"user":     conn.Config().User,
			"database": conn.Config().Database,
			"sql":      data.SQL,
			"args":     data.Args,
		},
	).Debug("query execution start")
	return ctx
}

func NewLogrusQueryTracer(logger *logrus.Logger) pgx.QueryTracer {
	return &queryTracer{
		log: logger.WithField("layer", "infrastructure.postgres.queryTracer"),
	}
}
