package voicevox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	APIKey  string
}

func New(base string, apiKey string) *Client {
	return &Client{
		BaseURL: base,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: apiKey,
	}
}

func (c *Client) Synthesize(
	ctx context.Context,
	text string,
	speakerID string,
	speakerSpeed float64,
) ([]byte, error) {

	// ---- audio_query ----
	q := url.Values{}
	q.Add("text", text)
	q.Add("speaker", speakerID)

	req1, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/audio_query?"+q.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if c.APIKey != "" {
		req1.Header.Set("Authorization", "ApiKey "+c.APIKey)
	}

	resp, err := c.HTTP.Do(req1)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("audio_query failed: %s", string(b))
	}

	queryBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// ---- speed 変更 ----
	var query map[string]any
	if err := json.Unmarshal(queryBody, &query); err != nil {
		return nil, err
	}

	query["speedScale"] = speakerSpeed

	modified, _ := json.Marshal(query)

	// ---- synthesis ----
	req2, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/synthesis?speaker="+speakerID,
		bytes.NewReader(modified),
	)
	if err != nil {
		return nil, err
	}
	req2.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req2.Header.Set("Authorization", "ApiKey "+c.APIKey)
	}

	res2, err := c.HTTP.Do(req2)
	if err != nil {
		return nil, err
	}
	defer res2.Body.Close()

	if res2.StatusCode != 200 {
		b, _ := io.ReadAll(res2.Body)
		return nil, fmt.Errorf("synthesis failed: %s", string(b))
	}

	return io.ReadAll(res2.Body)
}

// GetSpeakers はVOICEVOXの話者一覧を取得します
func (c *Client) GetSpeakers(ctx context.Context) ([]Speaker, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.BaseURL+"/speakers",
		nil,
	)
	if err != nil {
		return nil, err
	}

	if c.APIKey != "" {
		req.Header.Set("Authorization", "ApiKey "+c.APIKey)
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("speakers failed: %s", string(b))
	}

	var speakers []Speaker
	if err := json.NewDecoder(res.Body).Decode(&speakers); err != nil {
		return nil, err
	}

	return speakers, nil
}
