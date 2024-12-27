package clients

import (
	"encoding/json"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"io"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoyaltyClient struct {
	client  httpClient
	baseUrl string
}

func NewLoyaltyClient(client httpClient, baseUrl string) *LoyaltyClient {
	return &LoyaltyClient{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (l *LoyaltyClient) GetLoyaltyByUser(username string) (models.LoyaltyInfoResponse, error) {
	urlReq := fmt.Sprintf("%s/%s", l.baseUrl, "loyalty")
	req, err := http.NewRequest(http.MethodGet, urlReq, nil)
	if err != nil {
		return models.LoyaltyInfoResponse{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("X-User-Name", username)
	resp, err := l.client.Do(req)
	if err != nil {
		return models.LoyaltyInfoResponse{}, fmt.Errorf("failed to make request: %w", err)
	}

	if resp == nil {
		return models.LoyaltyInfoResponse{}, models.EmptyResponseError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.LoyaltyInfoResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var respModel models.LoyaltyInfoResponse
		if err := json.Unmarshal(body, &respModel); err != nil {
			return models.LoyaltyInfoResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return respModel, nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		var respErr models.ErrorResponse
		if err := json.Unmarshal(body, &respErr); err != nil {
			return models.LoyaltyInfoResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		respErr.StatusCode = resp.StatusCode
		return models.LoyaltyInfoResponse{}, respErr
	default:
		return models.LoyaltyInfoResponse{}, models.ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    models.UndefinedResponseCodeError.Error(),
		}
	}
}

func (l *LoyaltyClient) DecreaseLoyalty(username string) error {
	urlReq := fmt.Sprintf("%s/%s", l.baseUrl, "reservations/decrease")
	req, err := http.NewRequest(http.MethodDelete, urlReq, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("X-User-Name", username)
	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	if resp == nil {
		return models.EmptyResponseError
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError:
		var respErr models.ErrorResponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if err := json.Unmarshal(body, &respErr); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		resp.Body.Close()
		respErr.StatusCode = resp.StatusCode
		return respErr
	default:
		return models.ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    models.UndefinedResponseCodeError.Error(),
		}
	}
}

func (l *LoyaltyClient) IncreaseLoyalty(username string) error {
	urlReq := fmt.Sprintf("%s/%s", l.baseUrl, "reservations/increase")
	req, err := http.NewRequest(http.MethodPost, urlReq, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("X-User-Name", username)
	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	if resp == nil {
		return models.EmptyResponseError
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError:
		var respErr models.ErrorResponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if err := json.Unmarshal(body, &respErr); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		resp.Body.Close()
		respErr.StatusCode = resp.StatusCode
		return respErr
	default:
		return models.ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    models.UndefinedResponseCodeError.Error(),
		}
	}
}
