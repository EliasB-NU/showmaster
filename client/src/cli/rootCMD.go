package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCMD = &cobra.Command{
	Use:   "showmaster-client",
	Short: "Showmaster Client",
	Long:  "Showmaster Client, connect to your showmaster application to run midi or gpio actions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the showmaster client")
		fmt.Println("For help run 'showmaster-client help'")
	},
}

func Execute() {
	rootCMD.AddCommand(helpCommand)
	rootCMD.AddCommand(midiCommand)
	rootCMD.AddCommand(gpioCommand)

	if err := rootCMD.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
