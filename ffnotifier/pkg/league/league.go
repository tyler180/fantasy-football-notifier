package league

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// var leagueIDRegex = regexp.MustCompile(`"id":"(\d+)"`)

const (
	leagueID = "LEAGUE_ID"
	year     = "2024"
	proto    = "https"
	apiHost  = "api.myfantasyleague.com"
	// json     = 1
)

type Leagues struct {
	Leagues LeagueContainer `json:"leagues"`
}

type LeagueContainer struct {
	League []League `json:"league"`
}

type League struct {
	LeagueID    string `json:"league_id"`
	Name        string `json:"name"`
	FranchiseID string `json:"franchise_id"`
	URL         string `json:"url"`
}

func GetLeagueInfo(cookie string) ([]League, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s://%s/%s/export", proto, apiHost, year)
	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("MFL_USER_ID=%s", cookie))
	args := fmt.Sprintf("TYPE=myleagues&YEAR=%s&JSON=1", year)
	mlURL := fmt.Sprintf("%s?%s", url, args)

	req, err := http.NewRequest("GET", mlURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header = headers

	mlResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making league request: %v", err)
	}
	defer mlResp.Body.Close()

	mlBody, err := io.ReadAll(mlResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading league response: %v", err)
	}

	var leaguesResp Leagues
	err = json.Unmarshal(mlBody, &leaguesResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling league response: %v", err)
	}

	return leaguesResp.Leagues.League, nil
}
