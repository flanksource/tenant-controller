package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/tenant-controller/pkg/config"
	"github.com/flanksource/tenant-controller/pkg/tenant"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echopprof "github.com/sevennt/echo-pprof"
	"github.com/spf13/cobra"
)

var debug bool
var configFile string

var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start a server and accept API requests",
	Run: func(_ *cobra.Command, _ []string) {
		serve(configFile)
	},
}

func serve(configFile string) {
	if configFile == "" {
		log.Fatalln("Must specify the config file")
	}

	if err := config.SetConfig(configFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedCors,
	}))
	e.Use(middleware.Logger())

	if debug {
		logger.Infof("Starting pprof at /debug")
		echopprof.Wrap(e)
	}

	var err error
	tenant.ClerkTenantWebhook, err = tenant.NewWebhook(config.Config.Clerk.WebhookSecret)
	if err != nil {
		log.Fatalf("Error setting up webhook: %v", err)
	}

	e.GET("/health", func(c echo.Context) error { return c.JSON(200, map[string]string{"message": "ok"}) })
	e.POST("/tenant", tenant.CreateTenant)

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", httpPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Infof("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

func init() {
	ServerFlags(Serve.Flags())
	debugDefault := os.Getenv("DEBUG") == "true"
	Serve.Flags().BoolVar(&debug, "debug", debugDefault, "If true, start pprof at /debug")
	Serve.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
}
