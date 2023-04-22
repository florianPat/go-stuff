package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var verbose bool

var helloWorldCmd = &cobra.Command{
	Use: "hello",
	Short: "Just a hello world!",
	Args: cobra.MinimumNArgs(1),
	//Args: func(cmd *cobra.Command, args []string) error {
	//	if len(args) < 1 {
	//		return errors.New("requires a color argument")
	//	}
	//	if myapp.IsValidColor(args[0]) {
	//		return nil
	//	}
	//	return fmt.Errorf("invalid color specified: %s", args[0])
	//},
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.Println("Verbose output enabled!")
		}
		log.Println("Hello World!")
		log.Printf("%v", args[0])
	},
}

func init() {
	helloWorldCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Do you want verbosity?")
	rootCmd.AddCommand(helloWorldCmd)
}
