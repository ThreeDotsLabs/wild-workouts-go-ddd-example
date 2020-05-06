package main

import (
	"net/http"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server/httperr"
	"github.com/go-chi/render"
)

type HttpServer struct {
	db db
}

func (h HttpServer) GetTrainerAvailableHours(w http.ResponseWriter, r *http.Request) {
	queryParams := r.Context().Value("GetTrainerAvailableHoursParams").(*GetTrainerAvailableHoursParams)

	if queryParams.DateFrom.After(queryParams.DateTo) {
		httperr.BadRequest("date-from-after-date-to", nil, w, r)
		return
	}

	dates, err := h.db.GetDates(r.Context(), queryParams)
	if err != nil {
		httperr.InternalError("unable-to-get-dates", err, w, r)
		return
	}

	render.Respond(w, r, dates)
}

func (h HttpServer) MakeHourAvailable(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorised("no-user-found", err, w, r)
		return
	}

	if user.Role != "trainer" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	if err := h.db.UpdateAvailability(r, true); err != nil {
		httperr.InternalError("unable-to-update-availability", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) MakeHourUnavailable(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorised("no-user-found", err, w, r)
		return
	}
	if user.Role != "trainer" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	if err := h.db.UpdateAvailability(r, false); err != nil {
		httperr.InternalError("unable-to-update-availability", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
