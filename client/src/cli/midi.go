package cli

import "github.com/spf13/cobra"

var midiCommand = &cobra.Command{
	Use:   "midi",
	Short: "MIDI Command Line Interface",
	Long:  "MIDI Command Line Interface",
	Run:   midiCommandExecute,
}

func midiCommandExecute(cmd *cobra.Command, args []string) {

}
