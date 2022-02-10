package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ErrNonSuccessfulResponse = errors.New("server returned a non-200 response")

type StockLevelClient struct {
	client  *http.Client
	contact string
}

func New(client *http.Client, contact string) *StockLevelClient {
	return &StockLevelClient{
		client:  client,
		contact: contact,
	}
}

func (s *StockLevelClient) StockLevels(ctx context.Context, country, artNo string) (*Response, error) {
	stockLevelURL := fmt.Sprintf(
		"https://api.ingka.ikea.com/cia/availabilities/ru/%s?itemNos=%s&expand=StoresList",
		strings.ToLower(country),
		artNo,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, stockLevelURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json;version=2")
	req.Header.Set("User-Agent", fmt.Sprintf("github.com/patrick246/blahaj-exporter, contact: %s", s.contact))
	req.Header.Set("X-Client-ID", "b6c117e5-ae61-4ef5-b4cc-e0b1e37f0631")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status=%d, response=%s", ErrNonSuccessfulResponse, res.StatusCode, string(responseBytes))
	}

	var response Response
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
