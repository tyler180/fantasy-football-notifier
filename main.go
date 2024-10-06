package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const (
	leagueID = "79286"
	// username = "USERNAME"
	// password = "PASSWORD"
	year    = "2024"
	proto   = "https"
	apiHost = "api.myfantasyleague.com"
	json    = 0
	reqType = "league"
)

func main() {
	client := &http.Client{}

	username, password, err := LoadCredentials()
	if err != nil {
		fmt.Printf("Error loading credentials: %v\n", err)
		return
	}

	loginURL := fmt.Sprintf("https://%s/%s/login?USERNAME=%s&PASSWORD=%s&XML=1", apiHost, year, username, password)
	fmt.Printf("Making request to get cookie: %s\n", loginURL)
	loginResp, err := client.Get(loginURL)
	if err != nil {
		fmt.Printf("Error making login request: %v\n", err)
		return
	}
	defer loginResp.Body.Close()

	body, err := ioutil.ReadAll(loginResp.Body)
	if err != nil {
		fmt.Printf("Error reading login response: %v\n", err)
		return
	}

	cookieRegex := regexp.MustCompile(`MFL_USER_ID="([^"]*)">OK`)
	matches := cookieRegex.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		fmt.Printf("Cannot get login cookie. Response: %s\n", string(body))
		return
	}
	cookie := matches[1]
	fmt.Printf("Got cookie %s\n", cookie)

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

	mlBody, err := ioutil.ReadAll(mlResp.Body)
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
}

func LoadCredentials() (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", fmt.Errorf("error loading .env file: %v", err)
	}
	godotenv.Load()

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	return username, password, nil
}
