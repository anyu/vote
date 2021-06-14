package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flights",
	Long:  "List all flights",
	RunE:  list,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(c *cobra.Command, args []string) error {
	err := listFlights()
	if err != nil {
		return err
	}
	return nil
}

func listFlights() error {
	apiKey := os.Getenv("RAPID_API_KEY")
	apiHost := os.Getenv("RAPID_API_HOST")
	country := "us"
	currency := "usd"
	locale := "en-US"
	destination := "lond-sky"
	origin := "pari-sky"
	outDate := "anytime"
	inDate := "anytime"

	url := fmt.Sprintf("https://skyscanner-skyscanner-flight-search-v1.p.rapidapi.com/apiservices/browsequotes/v1.0/%s/%s/%s/%s/%s/%s?inboundpartialdate=%s",
		country, currency, locale, origin, destination, outDate, inDate)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}
