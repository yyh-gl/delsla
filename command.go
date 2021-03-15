package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	perform bool
	// FIXME: variable name
	days int

	rootCmd = &cobra.Command{
		Use:   "delsla",
		Short: "Delsla is the tool for deleting slack messages",
		Long:  "Delsla is the tool for deleting slack messages.",
		Run: func(cmd *cobra.Command, args []string) {
			chs, err := getChannels()
			if err != nil {
				log.Fatal(err)
			}

			deleteMessagesCount := 0
			for _, ch := range chs {
				mss, err := getMessages(ch.ID, days)
				if err != nil {
					log.Fatal(err)
				}

				if !perform || verbose {
					for _, m := range mss {
						fmt.Println(m.Text)
					}
				}

				if perform {
					if err := deleteMessages(ch.ID, mss); err != nil {
						log.Fatal(err)
					}
				}

				deleteMessagesCount += len(mss)
			}

			if perform {
				switch deleteMessagesCount {
				case 0:
					fmt.Println("\nNo deleted messages")
				case 1:
					fmt.Println("\nDelete 1 message")
				default:
					fmt.Printf("\nDelete %d messages\n", deleteMessagesCount)
				}
			} else {
				fmt.Println("\nThis is dry-run")
			}

			os.Exit(0)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  "Print the version number of Delsla",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Delsla v0.0.4")
			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&perform, "perform", "p", false, "perform")
	rootCmd.PersistentFlags().IntVarP(&days, "days", "d", 3, "delete messages older than {days} days")

	rootCmd.AddCommand(versionCmd)
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
