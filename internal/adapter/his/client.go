package his

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"encoding/json"
	"fmt"
	"net/http"
)

type hisClient struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) port.IHisClient {
	return &hisClient{
		httpClient: httpClient,
	}
}

const SEARCH_URL = "%s/patient/%s"

func (c *hisClient) FetchPatient(apiBase, id string) (*models.HISPatientResponse, error) {
	url := fmt.Sprintf(SEARCH_URL, apiBase, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HIS request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HIS returned status %d", resp.StatusCode)
	}

	var patient models.HISPatientResponse
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, fmt.Errorf("failed to decode HIS response: %w", err)
	}
	return &patient, nil
}
