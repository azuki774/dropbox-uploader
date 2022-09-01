package cmd

import (
	"azuki774/dropbox-uploader/internal/logger"
	"azuki774/dropbox-uploader/internal/uploader"
	"fmt"

	"github.com/spf13/cobra"
)

var opt uploader.UploadOption

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpload(&opt)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	fl := uploadCmd.Flags()
	fl.BoolVarP(&opt.Overwrite, "overwrite", "o", false, "overwrite if exists same name file")
	fl.BoolVarP(&opt.Update, "update", "u", false, "update when exists same name file if it has an newer timestump.")
	fl.StringVarP(&opt.SrcDir, "src-dir", "s", "", "source file or directory")
	fl.StringVarP(&opt.DstDir, "dst-dir", "d", "", "Dropbox target directory")
	fl.StringVarP(&opt.AccessToken, "token", "t", "", "Dropbox access token")
	fl.BoolVarP(&opt.Dryrun, "dry-run", "", false, "dry run")
}

func runUpload(opt *uploader.UploadOption) (err error) {
	opt.Logger, err = logger.NewLogger()
	if err != nil {
		fmt.Println("logger initialize error: %w", err)
		return err
	}
	defer opt.Logger.Sync()
	return uploader.Run(opt)
}
