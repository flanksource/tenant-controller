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
	"github.com/flanksource/tenant-controller/pkg"
	"github.com/flanksource/tenant-controller/pkg/api"
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
	Run: func(_ *cobra.Command, configFiles []string) {
		if configFile == "" {
			serve(configFiles[0])
		}
		serve(configFile)
	},
}

func serve(configFile string) {
	if len(configFile) == 0 {
		log.Fatalln("Must specify the config file")
	}

	if err := pkg.SetConfig(configFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedCors,
	}))

	if debug {
		logger.Infof("Starting pprof at /debug")
		echopprof.Wrap(e)
	}

	e.POST("/tenant", api.CreateTenant)

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