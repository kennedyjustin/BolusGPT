package dexcom

// Ported from https://github.com/gagebenne/pydexcom

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	ApplicationId           = "d89443d2-327c-4a6f-89e5-496bbb0317db"
	BaseUrl                 = "https://share2.dexcom.com/ShareWebServices/Services"
	AuthEndpoint            = "/General/AuthenticatePublisherAccount"
	LoginEndpoint           = "/General/LoginPublisherAccountById"
	GlucoseReadingsEndpoint = "/Publisher/ReadPublisherLatestGlucoseValues"
)

type Client struct {
	Username  string
	Password  string
	AccountId string
	SessionId string
}

type ClientInput struct {
	Username string
	Password string
}

func NewClient(input ClientInput) (*Client, error) {
	client := &Client{
		Username: input.Username,
		Password: input.Password,
	}

	err := client.RetrieveAccountId()
	if err != nil {
		return nil, err
	}

	err = client.RetrieveSessionId()
	if err != nil {
		return nil, err
	}

	return client, nil
}

type AuthRequest struct {
	AccountName   string `json:"accountName"`
	Password      string `json:"password"`
	ApplicationId string `json:"applicationId"`
}

func (c *Client) RetrieveAccountId() error {
	authRequest := AuthRequest{
		AccountName:   c.Username,
		Password:      c.Password,
		ApplicationId: ApplicationId,
	}

	requestBody, err := json.Marshal(&authRequest)
	if err != nil {
		return err
	}

	response, err := http.Post(BaseUrl+AuthEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	c.AccountId = strings.Trim(string(responseBody), "\"")

	err = uuid.Validate(c.AccountId)
	if err != nil {
		return err
	}

	return nil
}

type LoginRequest struct {
	AccountId     string `json:"accountId"`
	Password      string `json:"password"`
	ApplicationId string `json:"applicationId"`
}

func (c *Client) RetrieveSessionId() error {
	loginRequest := LoginRequest{
		AccountId:     c.AccountId,
		Password:      c.Password,
		ApplicationId: ApplicationId,
	}

	requestBody, err := json.Marshal(&loginRequest)
	if err != nil {
		return err
	}

	response, err := http.Post(BaseUrl+LoginEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	c.SessionId = strings.Trim(string(responseBody), "\"")
	err = uuid.Validate(c.SessionId)
	if err != nil {
		return err
	}

	return nil
}
