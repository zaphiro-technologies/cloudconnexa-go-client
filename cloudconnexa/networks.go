package cloudconnexa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NetworkConnector struct {
	Description     string `json:"description"`
	Id              string `json:"id"`
	IPv4Address     string `json:"ipV4Address"`
	IPv6Address     string `json:"ipV6Address"`
	Name            string `json:"name"`
	NetworkItemId   string `json:"networkItemId"`
	NetworkItemType string `json:"networkItemType"`
	VpnRegionId     string `json:"vpnRegionId"`
}

type Network struct {
	Connectors     []NetworkConnector `json:"connectors"`
	Description    string             `json:"description"`
	Egress         bool               `json:"egress"`
	Id             string             `json:"id"`
	InternetAccess string             `json:"internetAccess"`
	Name           string             `json:"name"`
	Routes         []Route            `json:"routes"`
	SystemSubnets  []string           `json:"systemSubnets"`
}

type NetworkPageResponse struct {
	Content          []Network `json:"content"`
	NumberOfElements int       `json:"numberOfElements"`
	Page             int       `json:"page"`
	Size             int       `json:"size"`
	Success          bool      `json:"success"`
	TotalElements    int       `json:"totalElements"`
	TotalPages       int       `json:"totalPages"`
}

type NetworksService service

func (c *NetworksService) GetByPage(page int, size int) (NetworkPageResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/networks/page?page=%d&size=%d", c.client.BaseURL, page, size), nil)
	if err != nil {
		return NetworkPageResponse{}, err
	}

	body, err := c.client.DoRequest(req)
	if err != nil {
		return NetworkPageResponse{}, err
	}

	var response NetworkPageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return NetworkPageResponse{}, err
	}

	return response, nil
}

func (c *NetworksService) List() ([]Network, error) {
	var allNetworks []Network
	pageSize := 10
	page := 0

	for {
		response, err := c.GetByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		allNetworks = append(allNetworks, response.Content...)

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return allNetworks, nil
}

func (c *NetworksService) GetByName(name string) (*Network, error) {
	networks, err := c.List()
	if err != nil {
		return nil, err
	}

	for _, n := range networks {
		if n.Name == name {
			return &n, nil
		}
	}
	return nil, nil
}

func (c *NetworksService) Get(networkId string) (*Network, error) {
	networks, err := c.List()
	if err != nil {
		return nil, err
	}

	for _, n := range networks {
		if n.Id == networkId {
			return &n, nil
		}
	}
	return nil, nil
}

func (c *NetworksService) Create(network Network) (*Network, error) {
	networkJson, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/networks", c.client.BaseURL), bytes.NewBuffer(networkJson))
	if err != nil {
		return nil, err
	}

	body, err := c.client.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var n Network
	err = json.Unmarshal(body, &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (c *NetworksService) Update(network Network) error {
	networkJson, err := json.Marshal(network)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/networks/%s", c.client.BaseURL, network.Id), bytes.NewBuffer(networkJson))
	if err != nil {
		return err
	}

	_, err = c.client.DoRequest(req)
	return err
}

func (c *NetworksService) Delete(networkId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/networks/%s", c.client.BaseURL, networkId), nil)
	if err != nil {
		return err
	}

	_, err = c.client.DoRequest(req)
	return err
}
