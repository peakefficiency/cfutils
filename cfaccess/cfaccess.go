package cfaccess

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudflare/cloudflared/carrier"
	"github.com/cloudflare/cloudflared/token"
	"github.com/pkg/errors"
)

// DefaultClientTimeout specifies the time limit for requests made by the HTTP client.
const DefaultClientTimeout = 5 * time.Second

// RequestWithCloudflareAccess performs an HTTP request to a Cloudflare Access protected URL.
// It automatically handles Access token retrieval and usage.
// `customHeaders` is optional and can be nil if no additional headers are needed.
func GetWithAccess(appURL string, customHeaders ...http.Header) ([]byte, error) {
	parsedURL, err := url.Parse(appURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL format")
	}

	appInfo, err := token.GetAppInfo(parsedURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get app info")
	}

	// Attempt to get an existing token, if it's not valid, fetch a new one
	tok, err := token.GetAppTokenIfExists(appInfo)
	if err != nil || tok == "" {
		tok, err = token.FetchToken(parsedURL, appInfo, nil) // Pass nil as logger since we're not using Sentry
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch token")
		}
	}

	// Create the HTTP client and request
	client := &http.Client{
		Timeout: DefaultClientTimeout,
	}
	req, err := http.NewRequest("GET", appURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set the Access token header
	req.Header.Set(carrier.CFAccessTokenHeader, tok)

	// Add any custom headers if provided
	for _, headers := range customHeaders {
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return body, nil
}
