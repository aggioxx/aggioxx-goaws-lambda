package http

import (
	"awesomeProject1/internal/core/domain/model"
	"awesomeProject1/internal/core/port"
	"awesomeProject1/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type EventHandler struct {
	EventService port.EventInboundPort
	log          *logger.Logger
}

type EventosOperacaoPayload struct {
	IdPessoaFisica       string         `json:"idPessoaFisica"`
	IdPessoaJuridica     string         `json:"idPessoaJuridica"`
	NumeroContrato       string         `json:"numeroContrato"`
	IsMultiplosContratos bool           `json:"isMultiplosContratos"`
	IdTipoContrato       string         `json:"idTipoContrato"`
	Metadata             model.Metadata `json:"metadata"`
}

type EventResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewEventHandler(eventService port.EventInboundPort, log *logger.Logger) *EventHandler {
	return &EventHandler{
		EventService: eventService,
		log:          log,
	}
}

func (h *EventHandler) ProcessEventosOperacao(requestBody []byte) (EventResult, int) {
	var payload EventosOperacaoPayload
	if err := json.Unmarshal(requestBody, &payload); err != nil {
		h.log.Warnf("Invalid JSON body: %v", err)
		return EventResult{
			Success: false,
			Error:   fmt.Sprintf("Invalid JSON: %s", err.Error()),
		}, 400
	}

	eventosOperacao := &model.EventosOperacao{
		IdPessoaFisica:       payload.IdPessoaFisica,
		IdPessoaJuridica:     payload.IdPessoaJuridica,
		NumeroContrato:       payload.NumeroContrato,
		IsMultiplosContratos: payload.IsMultiplosContratos,
		IdTipoContrato:       payload.IdTipoContrato,
		Metadata:             payload.Metadata,
	}

	if !eventosOperacao.IsValid() {
		h.log.Warn("Missing required fields: at least one of idPessoaFisica or idPessoaJuridica must be provided")
		return EventResult{
			Success: false,
			Error:   "Missing required fields: at least one of idPessoaFisica or idPessoaJuridica must be provided",
		}, 400
	}

	h.log.Infof("Processing eventos_operacao for metadata.id: %s", eventosOperacao.Metadata.Id)
	if err := h.EventService.ProcessEventosOperacao(eventosOperacao); err != nil {
		h.log.Errorf("Failed to process event: %v", err)
		return EventResult{
			Success: false,
			Error:   err.Error(),
		}, 500
	}

	h.log.Info("Event processed successfully")
	return EventResult{
		Success: true,
		Message: "Event processed successfully",
	}, 200
}

func (h *EventHandler) HandleLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.log.Info("Processing Lambda request")

	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Access-Control-Allow-Methods": "POST,OPTIONS",
	}

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "POST" {
		h.log.Warnf("Method not allowed: %s", req.HTTPMethod)
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Headers:    headers,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}

	if !strings.HasSuffix(req.Path, "/eventos_operacao") {
		h.log.Warnf("Path not found: %s", req.Path)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    headers,
			Body:       `{"error": "Path not found"}`,
		}, nil
	}

	result, statusCode := h.ProcessEventosOperacao([]byte(req.Body))

	var responseBody string
	if result.Success {
		responseBody = `{"message": "` + result.Message + `"}`
	} else {
		responseBody = `{"error": "` + result.Error + `"}`
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       responseBody,
	}, nil
}
