package golatch1

import "fmt"
import "errors"

var (
	ErrAlreadyPaired = errors.New("LatchError 205 Account and application already paired")
	ErrTokenNotFound = errors.New("LatchError 206 Pairing token not found or expired")
)

const (
	statusOn  = "on"
	statusOff = "off"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type operationStatus struct {
	Status string `json:"status"`
}

type statusData struct {
	Operations map[string]operationStatus `json:"operations"`
}

type statusResponse struct {
	Data  *statusData    `json:"data,omitempty"`
	Error *errorResponse `json:"error,omitempty"`
}

type unpairResponse struct {
	Error *errorResponse `json:"error,omitempty"`
}

type pairingData struct {
	AccountId string `json:"accountId"`
}

type pairResponse struct {
	Data  *pairingData   `json:"data,omitempty"`
	Error *errorResponse `json:"error,omitempty"`
}

func latchError2Error(err *errorResponse) error {
	switch err.Code {
	case 205:
		return ErrAlreadyPaired
	case 206:
		return ErrTokenNotFound
	}
	return fmt.Errorf("LatchError %v %v", err.Code, err.Message)
}

func (app *LatchApp) Pair(token string) (accountId string, err error) {
	fullPath := "/api/1.0/pair/" + token
	var response pairResponse
	if err := app.call(httpGet, fullPath, nil, nil, &response); err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", latchError2Error(response.Error)
	}
	return response.Data.AccountId, nil
}

func (app *LatchApp) Unpair(accountId string) error {
	fullPath := "/api/1.0/unpair/" + accountId
	var response unpairResponse
	if err := app.call(httpGet, fullPath, nil, nil, &response); err != nil {
		return err
	}
	if response.Error != nil {
		return latchError2Error(response.Error)
	}
	return nil
}

func (app *LatchApp) StatusIsOn(accountId string) (bool, error) {
	fullPath := "/api/1.0/status/" + accountId
	var response statusResponse

	if err := app.call(httpGet, fullPath, nil, nil, &response); err != nil {
		return false, err
	}

	if response.Error != nil {
		return false, latchError2Error(response.Error)
	}

	if len(response.Data.Operations) != 1 {
		return false, fmt.Errorf("LatchProtocolError %v", len(response.Data.Operations))
	}

	for _, v := range response.Data.Operations {
		switch v.Status {
		case statusOn:
			return true, nil
		case statusOff:
			return false, nil
		}
	}

	return false, errors.New("LatchProtocolError invalid status")
}
