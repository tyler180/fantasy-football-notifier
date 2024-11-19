package players

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	leagueID = "LEAGUE_ID"
	year     = "2024"
	proto    = "https"
	apiHost  = "api.myfantasyleague.com"
	json     = 1
)

func FreeAgents(cookie, league_id, position string) error {
	client := &http.Client{}

	// cookie, err := cmd.GetCookie(client)
	// if err != nil {
	// 	return fmt.Errorf("error getting cookie: %v", err)
	// }

	url := fmt.Sprintf("%s://%s/%s/export", proto, apiHost, year)
	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("MFL_USER_ID=%s", cookie))
	args := fmt.Sprintf("TYPE=freeAgents&L=%s&W=&JSON=%d", league_id, json)
	if position != "" {
		args = fmt.Sprintf("TYPE=freeAgents&L=%s&W=&POS=%s&JSON=%d", league_id, position, json)
	}
	mlURL := fmt.Sprintf("%s?%s", url, args)

	req, err := http.NewRequest("GET", mlURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header = headers

	mlResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making league request: %v", err)
	}
	defer mlResp.Body.Close()

	mlBody, err := io.ReadAll(mlResp.Body)
	if err != nil {
		return fmt.Errorf("error reading league response: %v", err)
	}

	leagueHostRegex := regexp.MustCompile(`url="(https?)://([a-z0-9]+.myfantasyleague.com)/` + year + `/home/` + leagueID + `"`)
	leagueMatches := leagueHostRegex.FindStringSubmatch(string(mlBody))
	if len(leagueMatches) < 3 {
		fmt.Printf("In the players package. Cannot find league host in response: %s\n", string(mlBody))
		return nil
	}
	protocol := leagueMatches[1]
	leagueHost := leagueMatches[2]
	fmt.Printf("Got league host %s\n", leagueHost)
	url = fmt.Sprintf("%s://%s/%s/export", protocol, leagueHost, year)
	fmt.Println(url)

	// Ensure the program ends cleanly
	fmt.Println("Program completed successfully.")

	return nil
}
