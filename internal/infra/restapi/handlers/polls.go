package handlers

import (
	"github.com/babadro/forecaster/internal/infra/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type service interface{}

type Polls struct {
	svc service
}

func NewPolls(svc service) *Polls {
	return &Polls{svc: svc}
}

func (h *Polls) CreatePoll(params operations.CreatePollParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.CreatePoll has not yet been implemented")
}

func (h *Polls) GetPollByID(params operations.GetPollByIDParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.GetPoll has not yet been implemented")
}

func (h *Polls) UpdatePoll(params operations.UpdatePollParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.UpdatePoll has not yet been implemented")
}

func (h *Polls) DeletePoll(params operations.DeletePollParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.DeletePoll has not yet been implemented")
}

func (h *Polls) CreateOption(params operations.CreateOptionParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.CreateOption has not yet been implemented")
}

func (h *Polls) UpdateOption(params operations.UpdateOptionParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.UpdateOption has not yet been implemented")
}

func (h *Polls) DeleteOption(params operations.DeleteOptionParams) middleware.Responder {
	return middleware.NotImplemented("operation operations.DeleteOption has not yet been implemented")
}
