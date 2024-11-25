package httpclient

import (
	"errors"
	"io"
	"net/http"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

func (h *HTTPClient) Get(url string) (body string, title string, err error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("non-200 HTTP status: " + resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	return string(bytes), extractTitle(string(bytes)), nil
}

func extractTitle(html string) string {
	startTag := "<title>"
	endTag := "</title>"
	start := indexOf(html, startTag)
	if start == -1 {
		return ""
	}
	start += len(startTag)
	end := indexOf(html, endTag)
	if end == -1 || end < start {
		return ""
	}
	return html[start:end]
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
