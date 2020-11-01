package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wmolicki/wd/deployed"
)

var rootCmd = &cobra.Command{
	Use:   "wd",
	Short: "wd is a cmdline tool for checking versions of deployed services",
	Run: func(cmd *cobra.Command, args []string) {

		r, err := os.Open("./resources/services_conf.yml")

		if err != nil {
			log.Fatalf("could not open conf file: %v", err)
		}

		config, err := deployed.LoadConfig(r)
		if err != nil {
			log.Fatalf("could not load config: %v", err)
		}

		env := config.Def.Env
		services := deployed.LoadServices(config, env)

		versions := deployed.FetchVersions(services)
		for i, version := range versions {
			fmt.Printf("service %d: %s\n", i, version)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
