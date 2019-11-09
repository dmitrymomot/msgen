package main

import (
	"time"
	"github.com/rs/zerolog"
	"github.com/nats-io/nats.go"
)

func setupNatsConnOptions(log zerolog.Logger, opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Info().Err(err).Msgf("nats client disconnected, will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Info().Msgf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatal().Err(err).Msgf("nats connection closed")
	}))
	return opts
}