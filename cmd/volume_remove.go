package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var volumeRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo volume rm VOLUME_NAME",
	Short:   "Remove a volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		volume, err := client.FindVolume(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s volume in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one volume with that value in your account")
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("volume", defaultYes, volume.Name) == true {

			_, err = client.DeleteVolume(volume.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": volume.ID, "Name": volume.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The volume called %s with ID %s was deleted\n", utility.Green(volume.Name), utility.Green(volume.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
