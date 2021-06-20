package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var e electionsResponse
	err =json.Unmarshal(body, &e)
	if err != nil {
		return err
	}

	for _, election := range e.Elections {
		fmt.Println(election)
	}

	return nil
}
