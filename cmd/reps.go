package cmd

import (
	"fmt"
	"github.com/anyu/vote/internal/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(repsCmd)
}

var repsCmd = &cobra.Command{
	Use:   "reps",
	Short: "List representatives",
	Long:  "List representatives",
	RunE:  reps,
}

func reps(c *cobra.Command, args []string) error {
	cl := client.New(os.Getenv("API_HOST"), os.Getenv("API_KEY"))

	address, err := solicitAddress()
	if err != nil {
		return err
	}

	locResp, err := cl.GetLocalReps(address)
	if err != nil {
		return err
	}
	displayReps(locResp)

	return nil
}

func displayReps(locResp *client.LocalRepsResponse) {
	for _, o := range locResp.Offices {
		fmt.Println()
		fmt.Println(o.Name)
		fmt.Printf("• %s", locResp.Officials[o.OfficialIndices[0]].Name)
		fmt.Println()
		fmt.Printf("• %s", locResp.Officials[o.OfficialIndices[0]].Party)
		if len(locResp.Officials[o.OfficialIndices[0]].Emails) > 0 {
			fmt.Println()
			fmt.Printf("• %s", locResp.Officials[o.OfficialIndices[0]].Emails[0])
		}
		fmt.Println()
	}
}
