package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bityacht-exchange-api-server/api"
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/cache/memory"
	"bityacht-exchange-api-server/internal/pkg/email"
	"bityacht-exchange-api-server/internal/pkg/exchange"
	"bityacht-exchange-api-server/internal/pkg/kyc"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/rbac"
	"bityacht-exchange-api-server/internal/pkg/receipt"
	_ "bityacht-exchange-api-server/internal/pkg/schedule"
	"bityacht-exchange-api-server/internal/pkg/sms"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"bityacht-exchange-api-server/internal/service"

	"github.com/spf13/cobra"
)

var tlsConfig = &tls.Config{
	MinVersion:               tls.VersionTLS12,
	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	PreferServerCipherSuites: true,
	CipherSuites: []uint16{
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_CHACHA20_POLY1305_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	},
}

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Server will run the server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootCtx, rootCancel := context.WithCancel(context.Background())
		rbac.Init()
		exchange.Init()
		memory.Init(rootCtx)
		sms.Init()
		email.Init()
		receipt.MustInit()
		kyc.Init()
		wallet.MustInit()

		// Start services
		go startServices(rootCtx)

		// Start http server
		httpServer := &http.Server{
			Addr:              net.JoinHostPort(configs.Config.Server.Host, strconv.FormatUint(uint64(configs.Config.Server.Port), 10)),
			Handler:           api.NewRouter(),
			ReadTimeout:       configs.Config.Server.ReadTimeout,
			WriteTimeout:      configs.Config.Server.WriteTimeout,
			MaxHeaderBytes:    1 << 20, // 1 MB
			ReadHeaderTimeout: configs.Config.Server.ReadTimeout,
		}
		go func() {
			var err error

			serveWithTLS := configs.Config.Server.CertFile != "" && configs.Config.Server.KeyFile != ""
			logger.Logger.Info().Str("Addr", httpServer.Addr).Bool("TLS", serveWithTLS).Msg("HTTP Server Start Listening.")

			if serveWithTLS {
				httpServer.TLSConfig = tlsConfig
				err = httpServer.ListenAndServeTLS(configs.Config.Server.CertFile, configs.Config.Server.KeyFile)
			} else {
				err = httpServer.ListenAndServe()
			}

			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Logger.Err(err).Str("Addr", httpServer.Addr).Msg("ListenAndServe error.")
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		receivedSignal := <-quit
		rootCancel()
		logger.Logger.Info().Str("Signal", receivedSignal.String()).Msg("Shutting down the server by user.")

		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()

		if err := httpServer.Shutdown(closeCtx); err != nil {
			logger.Logger.Err(err).Msg("Shut down server error.")
		}

		return nil
	},
}

func startServices(ctx context.Context) {
	service.RestoreIssuingReceipts(ctx)
}
