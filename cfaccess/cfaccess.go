package cfaccess

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudflare/cloudflared/logger"
	"github.com/cloudflare/cloudflared/token"
)

// DefaultClientTimeout specifies the time limit for requests made by the HTTP client.
const DefaultClientTimeout = 5 * time.Second

func GetWithAccess(appURL string) ([]byte, error) {

	client := &http.Client{Timeout: DefaultClientTimeout}
	req, err := http.NewRequest("GET", appURL, nil)
	if err != nil {
		return nil, err
	}
	AccessURL, _ := url.Parse(appURL)
	cloudflaredLogger := logger.Create(nil)
	appInfo, err := token.GetAppInfo(AccessURL)
	if err != nil {
		return nil, err
	}
	token, err := token.FetchToken(AccessURL, appInfo, cloudflaredLogger)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cf-access-token", token)
	req.Header.Set("User-Agent", "github.com/peakefficiency/cfutils")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
