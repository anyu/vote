package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"strings"

	"github.com/anyu/vote/internal/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(electionsCmd)
}

var electionsCmd = &cobra.Command{
	Use:   "elections",
	Short: "List upcoming elections",
	Long:  "List upcoming elections",
	RunE:  elections,
}

func elections(c *cobra.Command, args []string) error {
	cl := client.New(os.Getenv("API_HOST"), os.Getenv("API_KEY"))

	eResp, err := cl.GetUpcomingElections()
	if err != nil {
		return err
	}

	chosenElection, err := solicitElectionInput(eResp)
	if err != nil {
		return err
	}

	prompt := promptui.Prompt{
		Label: "What's your address?",
	}

	address, err := prompt.Run()
	if err != nil {
		fmt.Printf("error running prompt: %v", err)
		return err
	}
	vResp, err := cl.GetVoterInfo(chosenElection.ID, address)
	if err != nil {
		fmt.Printf("error getting voter info: %v", err)
		return err
	}
	displayVotingInfo(vResp)

	return nil
}

func solicitElectionInput(eResp *client.ElectionsResponse) (*client.Election, error) {
	template := promptui.SelectTemplates{
		Active:   `❯️  {{ .Name }} |  {{ .ElectionDay.String }}`,
		Inactive: `{{ .Name }} |  {{ .ElectionDay.String }}`,
		Selected: `❯️️  {{ .Name }} |  {{ .ElectionDay.String }}`,
	}

	prompt := promptui.Select{
		Label:     "Pick an upcoming election",
		Items:     eResp.Elections,
		Templates: &template,
		HideHelp:  true,
	}
	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("running prompt failed: %v", err)
		return nil, err
	}
	return &eResp.Elections[i], nil
}

func displayVotingInfo(vResp *client.VoterInfoResponse) {
	for _, evs := range vResp.EarlyVoteSites {
		fmt.Println()
		fmt.Println(strings.ToUpper("Your early vote site:"))
		fmt.Println(evs.LocationName)
		fmt.Printf("%s\n%s, %s, %s\n", evs.Address.Line1, evs.Address.City, evs.Address.State, evs.Address.Zip)
		fmt.Println()
		fmt.Println(strings.ToUpper("Polling hours:"))
		fmt.Println(evs.PollingHours)
	}

	for _, pl := range vResp.PollingLocations {
		fmt.Println()
		fmt.Println(strings.ToUpper("Your vote day polling site:"))
		fmt.Println(pl.LocationName)
		fmt.Printf("%s\n%s, %s, %s\n", pl.Address.Line1, pl.Address.City, pl.Address.State, pl.Address.Zip)
		fmt.Println()
		fmt.Println(strings.ToUpper("Polling hours:"))
		fmt.Println(pl.PollingHours)
	}
	for _, c := range vResp.Contests {
		fmt.Println()
		fmt.Printf("Type: %s", c.Type)
		fmt.Println()
		fmt.Printf("Office: %s", c.Office)
		fmt.Println()
		fmt.Printf("District: %s", c.District.Name)
		fmt.Println()
		for _, can := range c.Candidates {
			fmt.Println()
			fmt.Println(can.Name)
			fmt.Printf("Party: %s\n", can.Party)
		}
		fmt.Println()
		fmt.Println("=======================")
	}
}
