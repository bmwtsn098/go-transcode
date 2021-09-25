package transcode

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/m1k1o/go-transcode/internal/api"
	"github.com/m1k1o/go-transcode/internal/config"
	"github.com/m1k1o/go-transcode/internal/http"
)

var Service *Main

func init() {
	Service = &Main{
		RootConfig:   &config.Root{},
		ServerConfig: &config.Server{},
	}
}

type Main struct {
	RootConfig   *config.Root
	ServerConfig *config.Server

	logger     zerolog.Logger
	apiManager *api.ApiManagerCtx
	server     *http.ServerCtx
}

func (main *Main) Preflight() {
	main.logger = log.With().Str("service", "main").Logger()
}

func (main *Main) Start() {
	config := main.ServerConfig

	main.apiManager = api.New(config)
	main.apiManager.Start()

	main.server = http.New(config)
	main.server.Mount(main.apiManager.Mount)
	main.server.Start()

	main.logger.Info().Msgf("serving streams from basedir %s: %s", config.BaseDir, config.Streams)
}

func (main *Main) Shutdown() {
	var err error

	err = main.server.Shutdown()
	main.logger.Err(err).Msg("server shutdown")

	err = main.apiManager.Shutdown()
	main.logger.Err(err).Msg("api shutdown")
}

func (main *Main) ServeCommand(cmd *cobra.Command, args []string) {
	main.logger.Info().Msg("starting main server")
	main.Start()
	main.logger.Info().Msg("main ready")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	main.logger.Warn().Msgf("received %s, attempting graceful shutdown", sig)
	main.Shutdown()
	main.logger.Info().Msg("shutdown complete")
}
