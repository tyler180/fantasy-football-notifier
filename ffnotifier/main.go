package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

const (
	// leagueID = "79286"
	// username = "USERNAME"
	// password = "PASSWORD"
	year    = "2024"
	proto   = "https"
	apiHost = "api.myfantasyleague.com"
	json    = 0
	reqType = "league"
)

var leagueID string
var username string
var password string

func main() {
	client := &http.Client{}

	leagueID = os.Getenv("LEAGUE_ID")
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	// https://www49.myfantasyleague.com/2024/export?TYPE=playerScores&L=79286&APIKEY=&W=AVG&YEAR=2024&PLAYERS=&POSITION=&STATUS=freeagent&RULES=&COUNT=50&JSON=1
	cookie, err := GetCookie(client)
	if err != nil {
		fmt.Printf("Error getting cookie: %v\n", err)
		return
	}

	url := fmt.Sprintf("%s://%s/%s/export", proto, apiHost, year)
	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("MFL_USER_ID=%s", cookie))
	mlArgs := fmt.Sprintf("TYPE=myleagues&JSON=%d", json)
	mlURL := fmt.Sprintf("%s?%s", url, mlArgs)
	fmt.Printf("Making request to get league host: %s\n", mlURL)

	req, err := http.NewRequest("GET", mlURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header = headers

	mlResp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making league request: %v\n", err)
		return
	}
	defer mlResp.Body.Close()

	mlBody, err := io.ReadAll(mlResp.Body)
	if err != nil {
		fmt.Printf("Error reading league response: %v\n", err)
		return
	}

	leagueHostRegex := regexp.MustCompile(`url="(https?)://([a-z0-9]+.myfantasyleague.com)/` + year + `/home/` + leagueID + `"`)
	leagueMatches := leagueHostRegex.FindStringSubmatch(string(mlBody))
	if len(leagueMatches) < 3 {
		fmt.Printf("Cannot find league host in response: %s\n", string(mlBody))
		return
	}
	newProto := leagueMatches[1]
	leagueHost := leagueMatches[2]
	fmt.Printf("Got league host %s\n", leagueHost)
	url = fmt.Sprintf("%s://%s/%s/export", newProto, leagueHost, year)
	fmt.Println(url)
}

func GetCookie(client *http.Client) (string, error) {
	loginURL := fmt.Sprintf("https://%s/%s/login?USERNAME=%s&PASSWORD=%s&XML=1", apiHost, year, username, password)
	fmt.Printf("Making request to get cookie: %s\n", loginURL)
	loginResp, err := client.Get(loginURL)
	if err != nil {
		return "", fmt.Errorf("error making login request: %v", err)
	}
	defer loginResp.Body.Close()

	body, err := io.ReadAll(loginResp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading login response: %v", err)
	}

	cookieRegex := regexp.MustCompile(`MFL_USER_ID="([^"]*)">OK`)
	matches := cookieRegex.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("cannot get login cookie. Response: %s", string(body))
	}
	cookie := matches[1]
	fmt.Printf("Got cookie %s\n", cookie)

	return cookie, nil
}
