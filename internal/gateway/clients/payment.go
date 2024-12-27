package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"io"
	"net/http"
)

type PaymentClient struct {
	client  httpClient
	baseUrl string
}

func NewPaymentClient(client httpClient, baseUrl string) *PaymentClient {
	return &PaymentClient{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (p *PaymentClient) GetByUUID(uuid string) (models.PaymentInfo, error) {
	urlReq := fmt.Sprintf("%s/%s/%s", p.baseUrl, "payments", uuid)
	req, err := http.NewRequest(http.MethodGet, urlReq, nil)
	if err != nil {
		return models.PaymentInfo{}, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return models.PaymentInfo{}, fmt.Errorf("failed to make request: %w", err)
	}

	if resp == nil {
		return models.PaymentInfo{}, models.EmptyResponseError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.PaymentInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var respModel models.PaymentInfo
		if err := json.Unmarshal(body, &respModel); err != nil {
			return models.PaymentInfo{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return respModel, nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError:
		var respErr models.ErrorResponse
		if err := json.Unmarshal(body, &respErr); err != nil {
			return models.PaymentInfo{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		respErr.StatusCode = resp.StatusCode
		return models.PaymentInfo{}, respErr
	default:
		return models.PaymentInfo{}, models.ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    models.UndefinedResponseCodeError.Error(),
		}
	}
}

func (p *PaymentClient) Cancel(uuid string) error {
	urlReq := fmt.Sprintf("%s/%s/%s", p.baseUrl, "reservations", uuid)
	req, err := http.NewRequest(http.MethodDelete, urlReq, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	if resp == nil {
		return models.EmptyResponseError
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest, http.StatusInternalServerError:
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

func (p *PaymentClient) CreatePayment(payment models.PaymentCreateRequest) (models.ExtendedPaymentInfo, error) {
	urlReq := fmt.Sprintf("%s/%s", p.baseUrl, "payments")
	reqBody, err := json.Marshal(payment)
	if err != nil {
		return models.ExtendedPaymentInfo{}, fmt.Errorf("failed to build request body: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, urlReq, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.ExtendedPaymentInfo{}, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return models.ExtendedPaymentInfo{}, fmt.Errorf("failed to make request: %w", err)
	}
	if resp == nil {
		return models.ExtendedPaymentInfo{}, models.EmptyResponseError
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ExtendedPaymentInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusCreated:
		var respModel models.ExtendedPaymentInfo
		if err := json.Unmarshal(body, &respModel); err != nil {
			return models.ExtendedPaymentInfo{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return respModel, nil
	case http.StatusBadRequest, http.StatusInternalServerError:
		var respErr models.ErrorResponse
		if err := json.Unmarshal(body, &respErr); err != nil {
			return models.ExtendedPaymentInfo{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		respErr.StatusCode = resp.StatusCode
		return models.ExtendedPaymentInfo{}, respErr
	default:
		return models.ExtendedPaymentInfo{}, models.ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    models.UndefinedResponseCodeError.Error(),
		}
	}
}
