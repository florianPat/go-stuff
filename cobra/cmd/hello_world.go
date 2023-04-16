package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "hello",
	Short: "Just a hello world!",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Hello World!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
