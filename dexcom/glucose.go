package dexcom

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type CurrentBloodGlucoseReadingRequest struct {
	SessionId string `json:"sessionId"`
	Minutes   int    `json:"minutes"`
	MaxCount  int    `json:"maxCount"`
}

type CurrentBloodGlucoseReadingError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type CurrentBloodGlucoseReading struct {
	Value int    `json:"Value"`
	Trend string `json:"Trend"`
}

// https://www.dexcom.com/all-access/dexcom-cgm-explained/trend-arrow-and-treatment-decisions
var TrendToDeltaMap = map[string]int{
	"DoubleUp":       45,
	"SingleUp":       30,
	"FortyFiveUp":    15,
	"Flat":           0,
	"FortyFiveDown":  -15,
	"SingleDown":     -30,
	"DoubleDown":     -45,
	"None":           0,
	"NotComputable":  0,
	"RateOutOfRange": 0,
}

func (c CurrentBloodGlucoseReading) Get15MinDeltaFromTrend() int {
	return TrendToDeltaMap[c.Trend]
}

func (c *Client) GetCurrentBloodGlucoseReading() (*CurrentBloodGlucoseReading, error) {
	reading, err := c.getCurrentBloodGlucoseReading()
	if err != nil {
		if strings.Contains(err.Error(), "retry session") {
			c.RetrieveSessionId()
			reading, err = c.getCurrentBloodGlucoseReading()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return reading, nil
}

func (c *Client) getCurrentBloodGlucoseReading() (*CurrentBloodGlucoseReading, error) {
	glucoseRequest := CurrentBloodGlucoseReadingRequest{
		SessionId: c.SessionId,
		Minutes:   10,
		MaxCount:  1,
	}

	requestBody, err := json.Marshal(&glucoseRequest)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(BaseUrl+GlucoseReadingsEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode > 300 {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var errorResponse CurrentBloodGlucoseReadingError
		err = json.Unmarshal(responseBody, &errorResponse)
		if err != nil {
			return nil, err
		}

		if errorResponse.Code == "SessionIdNotFound" || errorResponse.Code == "SessionNotValid" {
			return nil, errors.New("retry session")
		} else {
			return nil, errors.New(errorResponse.Code + ": " + errorResponse.Message)
		}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var glucoseResponse []*CurrentBloodGlucoseReading
	err = json.Unmarshal(responseBody, &glucoseResponse)
	if err != nil {
		return nil, err
	}

	return glucoseResponse[0], nil
}
