package cmd

import (
	"github.com/spf13/cobra"
)

var startOpt startOption

type startOption struct {
	SrcRootDir   string
	DstRootDir   string
	RefreshToken string
}

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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
