package cmd

import (
	"github.com/flanksource/commons/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var Root = &cobra.Command{
	Use: "tenant-controller",
}

var httpPort = 8080
var publicEndpoint = "http://localhost:8080"
var allowedCors []string

func ServerFlags(flags *pflag.FlagSet) {
	flags.IntVar(&httpPort, "httpPort", httpPort, "Port to expose the api endpoint on")
	flags.StringVar(&publicEndpoint, "public-endpoint", publicEndpoint, "Host on which the health dashboard is exposed. Could be used for generting-links, redirects etc.")
	flags.StringSliceVar(&allowedCors, "allowed-cors", []string{"*"}, "Allowed CORS origins")
}

func init() {
	logger.BindFlags(Root.PersistentFlags())
	Root.AddCommand(Serve)
}
