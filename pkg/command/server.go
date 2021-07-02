package command

import (
	"context"
	"strings"

	"github.com/cernbox/ocis-eosprojects/pkg/config"
	"github.com/cernbox/ocis-eosprojects/pkg/flagset"
	"github.com/cernbox/ocis-eosprojects/pkg/metrics"
	"github.com/cernbox/ocis-eosprojects/pkg/server/http"
	"github.com/cernbox/ocis-eosprojects/pkg/service/v0"
	svc "github.com/cernbox/ocis-eosprojects/pkg/service/v0"
	"github.com/cernbox/ocis-eosprojects/pkg/tracing"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/ocis-pkg/sync"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "start eosprojects service",
		Description: "List user EOS Projects",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			logger := NewLogger(cfg)
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// When running on single binary mode the before hook from the root command won't get called. We manually
			// call this before hook from ocis command, so the configuration can be loaded.
			if !cfg.Supervised {
				return ParseConfig(ctx, cfg)
			}
			logger.Debug().Str("service", "eosprojects").Msg("ignoring config file parsing when running supervised")
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			err := tracing.Configure(cfg, logger)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)
			mtrcs := metrics.New()

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			handler, err := svc.NewEosProjects(cfg.DB, svc.Logger(logger))
			if err != nil {
				logger.Error().Err(err).Msg("handler init")
				return err
			}

			handler = service.NewInstrument(handler, mtrcs)
			handler = service.NewLogging(handler, logger)
			handler = service.NewTracing(handler)

			httpServer := http.Server(
				http.Config(cfg),
				http.Logger(logger),
				http.Name(cfg.Server.Name),
				http.Context(ctx),
				http.Metrics(mtrcs),
				http.Handler(handler),
			)

			gr.Add(httpServer.Run, func(_ error) {
				logger.Info().Str("server", "http").Msg("shutting down server")
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// defineContext sets the context for the extension. If there is a context configured it will create a new child from it,
// if not, it will create a root context that can be cancelled.
func defineContext(cfg *config.Config) (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		if cfg.Context == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(cfg.Context)
	}()
}
