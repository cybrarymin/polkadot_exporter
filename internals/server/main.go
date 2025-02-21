package srv

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cybrarymin/polkadot_exporter/internals/collector"
	gsrpc "github.com/polkadot-go/api/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	BuildTime   string
	Version     string
	LogLevel    string
	ShowVersion bool
	ListenAddr  string
	CertPath    string
	CertKeyPath string
)

type Exporter struct {
	logger zerolog.Logger
	wg     sync.WaitGroup
	api    *gsrpc.SubstrateAPI
}

func NewExoprter(logger zerolog.Logger) (*Exporter, error) {
	api, err := gsrpc.NewSubstrateAPI(collector.RpcBackend)
	if err != nil {
		return &Exporter{
			logger,
			sync.WaitGroup{},
			api,
		}, err
	}
	return &Exporter{
		logger,
		sync.WaitGroup{},
		api,
	}, nil
}

func Start() {
	// Defining a new logger and setting condition to consider specific loglevel
	var nlogger zerolog.Logger
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if zerolog.LevelTraceValue == LogLevel {
		nlogger = zerolog.New(os.Stdout).With().Stack().Timestamp().Logger().Level(zerolog.TraceLevel)
	} else {
		loglvl, _ := zerolog.ParseLevel(LogLevel)
		nlogger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(loglvl)
	}

	exp, err := NewExoprter(nlogger)
	if err != nil {
		nlogger.Error().Err(err).Msg("couldn't establish connection to the rpc backend to scrape metrics. backend metrics won't be available")

	}

	collector.RegisterCollectors(&nlogger)

	scheme, host, err := ListenAddrParser(ListenAddr)
	if err != nil {
		exp.logger.Error().Err(err).Send()
		return
	}

	srv := &http.Server{
		Addr:         host,
		Handler:      exp.routes(),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(nlogger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownChan := make(chan error)
	go exp.gracefulShutdown(srv, shutdownChan)

	exp.logger.Info().Msgf("starting the exporter server on %s://%s.....", scheme, host)

	switch scheme {
	case "https":
		err = srv.ListenAndServeTLS(CertPath, CertKeyPath)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			exp.logger.Error().Err(err).Send()
			return
		}
	case "http":
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			exp.logger.Error().Err(err).Send()
			return
		}
	}

	err = <-shutdownChan // This channel will block main appliction not to finish until shutdown method return it's errors.
	if err != nil {
		exp.logger.Error().Err(err).Send()
	}

}

func (exp *Exporter) gracefulShutdown(srv *http.Server, shutdownErr chan error) {
	// Create a channel to redirect signal to it.
	quit := make(chan os.Signal, 1)

	// This will listen to signals specified and will relay to them to the channel specified.
	// This will impede program to exit by the signal
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Catching the signal and print it
	s := <-quit
	exp.logger.Info().Msgf("received os signal %s", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		shutdownErr <- err
	}

	exp.logger.Info().Msg("waiting for background tasks to finish")
	exp.wg.Wait()
	shutdownErr <- nil
	exp.logger.Info().Msg("stopped server")
}
