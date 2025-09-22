package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DTO koji Blog servis očekuje od Stakeholder servisa

type StakeholderClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewStakeholderClient(baseURL string) *StakeholderClient {
	return &StakeholderClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *StakeholderClient) GetUser(userID string) (*UserDTO, error) {
	url := fmt.Sprintf("%s/stakeholders/user/%s", c.BaseURL, userID) // ovo se slaže sa tvojim handlerom
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: %s", resp.Status)
	}

	var user UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
