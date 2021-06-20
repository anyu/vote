package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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
	ElectionDay string `json:"electionDay"`
	DivisionID  string `json:"ocdDivisionId"`
}

func getElections() error {
	apiHost := os.Getenv("API_HOST")
	apiKey := os.Getenv("API_KEY")
	ep := "elections"
	url := fmt.Sprintf("%s/%s?key=%s", apiHost, ep, apiKey)

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

	electionNameToID := make(map[string]string)
	var electionNames []string

	for _, e := range eResp.Elections {
		electionNameToID[e.Name] = e.ID
		electionDay, err := time.Parse(rawLayout, e.ElectionDay)
		if err != nil {
			fmt.Printf("error parsing election day: %s", err)
			return err
		}
		electionNames = append(electionNames, fmt.Sprintf("%s (%s)", e.Name, electionDay.Format(desiredLayout)))
	}
	prompt := promptui.Select{
		Label: "Pick an upcoming election",
		Items: electionNames,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("running prompt failed: %v", err)
		return err
	}
	fmt.Printf("You chose: %q", result)

	return nil
}
