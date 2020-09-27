package ports

import (
	"net/http"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server/httperr"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app/query"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/render"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (h HttpServer) GetTrainerAvailableHours(w http.ResponseWriter, r *http.Request) {
	queryParams := r.Context().Value("GetTrainerAvailableHoursParams").(*GetTrainerAvailableHoursParams)

	dateModels, err := h.app.Queries.TrainerAvailableHours.Handle(r.Context(), query.AvailableHours{
		From: queryParams.DateFrom,
		To:   queryParams.DateTo,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	dates := dateModelsToResponse(dateModels)
	render.Respond(w, r, dates)
}

func dateModelsToResponse(models []query.Date) []Date {
	var dates []Date
	for _, d := range models {
		var hours []Hour
		for _, h := range d.Hours {
			hours = append(hours, Hour{
				Available:            h.Available,
				HasTrainingScheduled: h.HasTrainingScheduled,
				Hour:                 h.Hour,
			})
		}

		dates = append(dates, Date{
			Date: openapi_types.Date{
				Time: d.Date,
			},
			HasFreeHours: d.HasFreeHours,
			Hours:        hours,
		})
	}

	return dates
}

func (h HttpServer) MakeHourAvailable(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "trainer" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	hourUpdate := &HourUpdate{}
	if err := render.Decode(r, hourUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.MakeHoursAvailable.Handle(r.Context(), hourUpdate.Hours)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) MakeHourUnavailable(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "trainer" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	hourUpdate := &HourUpdate{}
	if err := render.Decode(r, hourUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.MakeHoursUnavailable.Handle(r.Context(), hourUpdate.Hours)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
