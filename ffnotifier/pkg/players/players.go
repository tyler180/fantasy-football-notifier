package players

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/tyler180/fantasy-football-notifier/ffnotifier/cmd"
)

func freeAgents(ctx context.Context, league_id, position string) error {
	client := &http.Client{}

	cookie, err := cmd.GetCookie(client)
	if err != nil {
		return fmt.Errorf("Error getting cookie: %v", err)
	}

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
		return fmt.Errorf("Error creating request: %v", err)
	}
	req.Header = headers

	mlResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making league request: %v", err)
	}
	defer mlResp.Body.Close()

	mlBody, err := io.ReadAll(mlResp.Body)
	if err != nil {
		return fmt.Errorf("Error reading league response: %v", err)
	}

	return nil
}
