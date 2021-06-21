package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// go's magical reference date
	rawLayout     = "2006-01-02"
	desiredLayout = "January 2, 2006"
)

type client struct {
	apiHost string
	apiKey  string
}

func New(apiHost, apiKey string) *client {
	return &client{
		apiHost: apiHost,
		apiKey:  apiKey,
	}
}

type ElectionsResponse struct {
	Elections []Election `json:"elections"`
}

type Election struct {
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

type VoterInfoResponse struct {
	PollingLocations []Location `json:"pollingLocations"`
	EarlyVoteSites   []Location `json:"earlyVoteSites"`
	DropOffLocations []Location `json:"dropOffLocations"`
	State            []State    `json:"state"`
	Contests         []Contest  `json:"contests"`
}

type Location struct {
	Address      `json:"address"`
	PollingHours string  `json:"pollingHours"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	StartDate    string  `json:"startDate"`
	EndDate      string  `json:"endDate"`
}

type Address struct {
	LocationName string `json:"locationName"`
	Line1        string `json:"line1"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
}

type State struct {
	Name string `json:"name"`
}

type Contest struct {
	Type     string   `json:"type"`
	Office   string   `json:"office"`
	Level    []string `json:"level"`
	Roles    []string `json:"roles"`
	District struct {
		Name  string `json:"name"`
		Scope string `json:"scope"`
		ID    string `json:"id"`
	} `json:"district"`
	Candidates []struct {
		Name     string `json:"name"`
		Party    string `json:"party"`
		Channels []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"channels"`
	} `json:"candidates"`
}

func (c *client) GetUpcomingElections() (*ElectionsResponse, error) {
	var eResp *ElectionsResponse
	params := url.Values{}
	params.Add("key", c.apiKey)

	err := c.makeGetReq("elections", params, &eResp)
	if err != nil {
		return nil, err
	}

	return eResp, nil
}

func (c *client) GetVoterInfo(electionID, address string) (*VoterInfoResponse, error) {
	var vResp *VoterInfoResponse
	params := url.Values{}
	params.Add("electionId", electionID)
	params.Add("address", address)
	params.Add("key", c.apiKey)

	err := c.makeGetReq("voterinfo", params, &vResp)
	if err != nil {
		return nil, err
	}

	return vResp, nil
}

func (c *client) makeGetReq(ep string, params url.Values, respData interface{}) error {
	url := fmt.Sprintf("%s/%s?"+params.Encode(), c.apiHost, ep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	return nil
}
