package cmd

import (
	"azuki774/dropbox-uploader/internal/factory"
	"azuki774/dropbox-uploader/internal/logger"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

type StartOption struct {
	RepoInfo struct {
		Host string
		Port string
	}
}

var startOpt StartOption

// startCmd represents the upload command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := factory.NewClient()
		l, err := logger.NewLogger()
		if err != nil {
			fmt.Println(err)
			return err
		}

		tr := factory.NewTokenRepo(startOpt.RepoInfo.Host, startOpt.RepoInfo.Port)
		newTokenClient, err := factory.NewNewTokenClient(tr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		us := factory.NewUsecases(l, client, newTokenClient)
		srv := factory.NewServer(l, us)

		err = us.GetNewAccessToken()
		if err != nil {
			return err
		}

		c := cron.New()
		c.AddFunc("@hourly", func() { us.GetNewAccessToken() })
		c.Start()

		err = srv.Start()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&startOpt.RepoInfo.Host, "repo-host", "localhost", "token repository host")
	startCmd.Flags().StringVar(&startOpt.RepoInfo.Port, "repo-port", "80", "token repositroy port")
}
