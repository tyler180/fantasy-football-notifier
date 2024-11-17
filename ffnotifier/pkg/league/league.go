package league

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

func GetLeagueIDs(cookie string) ([][]string, error) {
	client := &http.Client{}

	// cookie, err := cmd.GetCookie(client)
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting cookie: %v", err)
	// }

	url := fmt.Sprintf("%s://%s/%s/export", proto, apiHost, year)
	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("MFL_USER_ID=%s", cookie))
	args := fmt.Sprintf("TYPE=myleagues&YEAR=%s&JSON=%d", year, json)
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

	leagueIDRegex := regexp.MustCompile(`"league_id":\s*"(\d+)"`)
	matches := leagueIDRegex.FindAllStringSubmatch(string(mlBody), -1)
	if matches == nil {
		fmt.Printf("No league_id found in response: %s\n", string(mlBody))
		return nil, nil
	}

	for _, match := range matches {
		if len(match) > 1 {
			fmt.Printf("Found league_id: %s\n", match[1])
			return matches, nil
		}
	}

	fmt.Println("Program completed successfully.")
	return nil, nil
}
