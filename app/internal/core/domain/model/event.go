package model

import "time"

type Event struct {
	ID          string                 `json:"id"`
	Payload     map[string]interface{} `json:"payload"`
	ProcessedAt string                 `json:"processed_at"`
}

type EventosOperacao struct {
	IdPessoaFisica       string    `json:"idPessoaFisica"`
	IdPessoaJuridica     string    `json:"idPessoaJuridica"`
	NumeroContrato       string    `json:"numeroContrato"`
	IsMultiplosContratos bool      `json:"isMultiplosContratos"`
	IdTipoContrato       string    `json:"idTipoContrato"`
	Metadata             Metadata  `json:"metadata"`
	EventType            string    `json:"eventType"`
	ProcessedAt          time.Time `json:"processedAt"`
}

type Metadata struct {
	Id              string `json:"id"`
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	DataCriacao     string `json:"dataCriacao"`
	DataAtualizacao string `json:"dataAtualizacao"`
}

func (e *EventosOperacao) IsValid() bool {
	return e.IdPessoaFisica != "" || e.IdPessoaJuridica != ""
}

func (e *EventosOperacao) DetermineEventType() string {
	if e.IdPessoaFisica != "" {
		return "pessoa_fisica_event"
	} else if e.IdPessoaJuridica != "" {
		return "pessoa_juridica_event"
	} else if e.NumeroContrato != "" {
		return "contrato_event"
	}
	return "generic_operacao_event"
}
