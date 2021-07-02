package command

import (
	"context"
	"os"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/cernbox/ocis-eosprojects/pkg/config"
	"github.com/cernbox/ocis-eosprojects/pkg/flagset"
	"github.com/cernbox/ocis-eosprojects/pkg/version"
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/spf13/viper"
)

// Execute is the entry point for the ocis-eosprojects command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "eosprojects",
		Version:  version.String,
		Usage:    "eosprojects, Extension to list user Eos Projects",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "CERNBox",
				Email: "cernbox-admins@cern.ch",
			},
		},

		Flags: flagset.RootWithConfig(cfg),

		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg)
		},

		Commands: []*cli.Command{
			Server(cfg),
			Health(cfg),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version,v",
		Usage: "Print the version",
	}

	return app.Run(os.Args)
}

// NewLogger initializes a service-specific logger instance.
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("eosprojects"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads hello configuration from Viper known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	sync.ParsingViperConfig.Lock()
	defer sync.ParsingViperConfig.Unlock()
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("EOSP")
	viper.AutomaticEnv()

	if c.IsSet("config-file") {
		viper.SetConfigFile(c.String("config-file"))
	} else {
		viper.SetConfigName("eosprojects")

		viper.AddConfigPath("/etc/ocis")
		viper.AddConfigPath("$HOME/.ocis")
		viper.AddConfigPath("./config")
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			logger.Info().
				Msg("Continue without config")
		case viper.UnsupportedConfigError:
			logger.Fatal().
				Err(err).
				Msg("Unsupported config type")
		default:
			logger.Fatal().
				Err(err).
				Msg("Failed to read config")
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatal().
			Err(err).
			Msg("Failed to parse config")
	}

	return nil
}

// SutureService allows for the hello command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
