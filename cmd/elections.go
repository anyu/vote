package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// go's magical reference date
	rawLayout     = "2006-01-02"
	desiredLayout = "January 2, 2006"
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

func init() {
	rootCmd.AddCommand(electionsCmd)
}

func elections(c *cobra.Command, args []string) error {
	err := getElections()
	if err != nil {
		return err
	}
	return nil
}

type electionsResponse struct {
	Elections []electionsData `json:"elections"`
}

type electionsData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ElectionDay eDay   `json:"electionDay"`
	DivisionID  string `json:"ocdDivisionId"`
}

type eDay string

func (e eDay) String() string {
	electionDay, _ := time.Parse(rawLayout, string(e))
	return electionDay.Format(desiredLayout)
}

func getElections() error {
	apiHost := os.Getenv("API_HOST")
	apiKey := os.Getenv("API_KEY")
	epElections := "elections"

	params := url.Values{}
	params.Add("key", apiKey)
	url := fmt.Sprintf("%s/%s?"+params.Encode(), apiHost, epElections)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var eResp electionsResponse
	err = json.NewDecoder(resp.Body).Decode(&eResp)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	template := promptui.SelectTemplates{
		Active:   `🗳️  {{ .Name }} |  {{ .ElectionDay.String }}`,
		Inactive: `{{ .Name }} |  {{ .ElectionDay.String }}`,
		Selected: `🗳  {{ .Name }} |  {{ .ElectionDay.String }}`,
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
		return err
	}
	chosenElection := eResp.Elections[i]
	fmt.Printf("You chose: %q", chosenElection.Name)

	prompt2 := promptui.Prompt{
		Label: "What's your address?",
	}

	address, err := prompt2.Run()
	if err != nil {
		fmt.Printf("error running prompt 2: %v", err)
		return err
	}

	err = getVoterInfo(chosenElection.ID, address)
	if err != nil {
		return err
	}

	return nil
}

type voterInfoResponse struct {
	PollingLocations []Location `json:"pollingLocations"`
	EarlyVoteSites   []Location `json:"earlyVoteSites"`
	DropOffLocations []Location `json:"dropOffLocations"`
	State            []State    `json:"state"`
}

type Location struct {
	Address      `json:"address"`
	PollingHours string  `json:"pollingHours"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	StartDate    string  `json:"startDate"`
	EndDate      string  `json:"endDate"`
}

type State struct {
	Name string `json:"name"`
}

type Address struct {
	LocationName string `json:"locationName"`
	Line1        string `json:"line1"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
}

func getVoterInfo(electionID, address string) error {
	apiHost := os.Getenv("API_HOST")
	apiKey := os.Getenv("API_KEY")

	epVoterInfo := "voterinfo"

	params := url.Values{}
	params.Add("electionId", electionID)
	params.Add("address", address)
	params.Add("key", apiKey)
	url := fmt.Sprintf("%s/%s?"+params.Encode(), apiHost, epVoterInfo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var vResp voterInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	for _, v := range vResp.EarlyVoteSites {
		fmt.Println(v.LocationName)
		fmt.Println(v.Address)
		fmt.Println(v.PollingHours)
	}

	return nil
}
