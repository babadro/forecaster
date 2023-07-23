package handlers

import (
	"context"
	"errors"

	"github.com/babadro/forecaster/internal/domain"
	"github.com/babadro/forecaster/internal/infra/restapi/operations"
	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
)

type service interface {
	GetSeriesByID(ctx context.Context, id int32) (models.Series, error)
	GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error)

	CreateSeries(ctx context.Context, s models.CreateSeries) (models.Series, error)
	CreatePoll(ctx context.Context, poll models.CreatePoll) (models.Poll, error)
	CreateOption(ctx context.Context, option models.CreateOption) (models.Option, error)

	UpdateSeries(ctx context.Context, id int32, s models.UpdateSeries) (models.Series, error)
	UpdatePoll(ctx context.Context, poll models.UpdatePoll) (models.Poll, error)
	UpdateOption(ctx context.Context, option models.UpdateOption) (models.Option, error)

	DeleteSeries(ctx context.Context, id int32) error
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, id int32) error
}

type Polls struct {
	svc service
}

func NewPolls(svc service) *Polls {
	return &Polls{svc: svc}
}

func (p *Polls) GetSeriesByID(params operations.GetSeriesByIDParams) middleware.Responder {
	series, err := p.svc.GetSeriesByID(params.HTTPRequest.Context(), params.SeriesID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return operations.NewGetSeriesByIDNotFound()
		}

		return operations.NewGetSeriesByIDInternalServerError()
	}

	return operations.NewGetSeriesByIDOK().WithPayload(&series)
}

func (p *Polls) GetPollByID(params operations.GetPollByIDParams) middleware.Responder {
	poll, err := p.svc.GetPollByID(params.HTTPRequest.Context(), params.PollID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return operations.NewGetPollByIDNotFound()
		}

		return operations.NewGetPollByIDInternalServerError()
	}
	return operations.NewGetPollByIDOK().WithPayload(&poll)
}

func (p *Polls) CreateSeries(params operations.CreateSeriesParams) middleware.Responder {
	s, err := p.svc.CreateSeries(params.HTTPRequest.Context(), *params.Series)
	if err != nil {
		return operations.NewCreateSeriesInternalServerError()
	}

	return operations.NewCreateSeriesCreated().WithPayload(&s)
}

func (p *Polls) CreatePoll(params operations.CreatePollParams) middleware.Responder {
	poll, err := p.svc.CreatePoll(params.HTTPRequest.Context(), *params.Poll)
	if err != nil {
		return operations.NewCreatePollInternalServerError()
	}

	return operations.NewCreatePollCreated().WithPayload(&poll)
}

func (p *Polls) CreateOption(params operations.CreateOptionParams) middleware.Responder {
	option, err := p.svc.CreateOption(params.HTTPRequest.Context(), *params.Option)
	if err != nil {
		return operations.NewCreateOptionInternalServerError()
	}

	return operations.NewCreateOptionCreated().WithPayload(&option)
}

func (p *Polls) UpdateSeries(params operations.UpdateSeriesParams) middleware.Responder {
	s, err := p.svc.UpdateSeries(params.HTTPRequest.Context(), params.SeriesID, *params.Series)
	if err != nil {
		return operations.NewUpdateSeriesInternalServerError()
	}

	return operations.NewUpdateSeriesOK().WithPayload(&s)
}

func (p *Polls) UpdatePoll(params operations.UpdatePollParams) middleware.Responder {
	poll, err := p.svc.UpdatePoll(params.HTTPRequest.Context(), *params.Poll)
	if err != nil {
		return operations.NewUpdatePollInternalServerError()
	}

	return operations.NewUpdatePollOK().WithPayload(&poll)
}

func (p *Polls) UpdateOption(params operations.UpdateOptionParams) middleware.Responder {
	option, err := p.svc.UpdateOption(params.HTTPRequest.Context(), *params.Option)
	if err != nil {
		return operations.NewUpdateOptionInternalServerError()
	}

	return operations.NewUpdateOptionOK().WithPayload(&option)
}

func (p *Polls) DeletePoll(params operations.DeletePollParams) middleware.Responder {
	err := p.svc.DeletePoll(params.HTTPRequest.Context(), params.PollID)
	if err != nil {
		return operations.NewDeletePollInternalServerError()
	}

	return operations.NewDeletePollNoContent()
}

func (p *Polls) DeleteOption(params operations.DeleteOptionParams) middleware.Responder {
	err := p.svc.DeleteOption(params.HTTPRequest.Context(), params.OptionID)
	if err != nil {
		return operations.NewDeleteOptionInternalServerError()
	}

	return operations.NewDeleteOptionNoContent()
}

func (p *Polls) DeleteSeries(params operations.DeleteSeriesParams) middleware.Responder {
	err := p.svc.DeleteSeries(params.HTTPRequest.Context(), params.SeriesID)
	if err != nil {
		return operations.NewDeleteSeriesInternalServerError()
	}

	return operations.NewDeleteSeriesNoContent()
}
