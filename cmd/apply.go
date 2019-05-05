package cmd

import (
	"io/ioutil"

	_ "github.com/moshloop/configadm/pkg"
	"github.com/moshloop/configadm/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Apply = cobra.Command{
		Use:   "apply",
		Short: "Apply the configuration to the local machine",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			cfg := GetConfig(cmd)

			files, commands, err := cfg.ApplyPhases()
			if err != nil {
				panic(err)
			}

			for path, file := range files {
				log.Infof("Writing %s\n", path)
				ioutil.WriteFile(path, []byte(file.Content), 0644)
			}

			for _, cmd := range commands {
				log.Infof("Executing %s\n", cmd.Cmd)
				if err := utils.Exec(cmd.Cmd); err != nil {
					log.Fatalf("Failed to run: %s, %s", cmd.Cmd, err)
				}
			}

		},
	}
)
