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

		conf, err := deployed.LoadConfig(r);
		if err != nil {
			log.Fatalf("could not load config file: %v", err)
		}

		services := deployed.LoadServices(conf, "qa")
		
		results := deployed.FetchVersions(services)
		for _, result := range results {
			fmt.Printf("%s: %s\n", result.Name, result.Version)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
