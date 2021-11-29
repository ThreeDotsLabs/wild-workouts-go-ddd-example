package ports

import (
	"context"
	"net/http"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server/httperr"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app/command"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app/query"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) GetTrainings(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var appTrainings []query.Training

	if user.Role == "trainer" {
		appTrainings, err = h.app.Queries.AllTrainings.Handle(r.Context())
	} else {
		appTrainings, err = h.app.Queries.TrainingsForUser.Handle(r.Context(), user)
	}

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	trainings := appTrainingsToResponse(appTrainings)
	trainingsResp := Trainings{trainings}

	render.Respond(w, r, trainingsResp)
}

func (h HttpServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	postTraining := PostTraining{}
	if err := render.Decode(r, &postTraining); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "attendee" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	cmd := command.ScheduleTraining{
		TrainingUUID: uuid.New().String(),
		UserUUID:     user.UUID,
		UserName:     user.DisplayName,
		TrainingTime: postTraining.Time,
		Notes:        postTraining.Notes,
	}
	err = h.app.Commands.ScheduleTraining.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "/trainings/"+cmd.TrainingUUID)
	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) CancelTraining(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.CancelTraining.Handle(r.Context(), command.CancelTraining{
		TrainingUUID: trainingUUID,
		User:         user,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	rescheduleTraining := PostTraining{}
	if err := render.Decode(r, &rescheduleTraining); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.RescheduleTraining.Handle(r.Context(), command.RescheduleTraining{
		User:         user,
		TrainingUUID: trainingUUID,
		NewTime:      rescheduleTraining.Time,
		NewNotes:     rescheduleTraining.Notes,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RequestRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	rescheduleTraining := PostTraining{}
	if err := render.Decode(r, &rescheduleTraining); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.RequestTrainingReschedule.Handle(r.Context(), command.RequestTrainingReschedule{
		User:         user,
		TrainingUUID: trainingUUID,
		NewTime:      rescheduleTraining.Time,
		NewNotes:     rescheduleTraining.Notes,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) ApproveRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.ApproveTrainingReschedule.Handle(r.Context(), command.ApproveTrainingReschedule{
		User:         user,
		TrainingUUID: trainingUUID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RejectRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.RejectTrainingReschedule.Handle(r.Context(), command.RejectTrainingReschedule{
		User:         user,
		TrainingUUID: trainingUUID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func appTrainingsToResponse(appTrainings []query.Training) []Training {
	var trainings []Training
	for _, tm := range appTrainings {
		t := Training{
			CanBeCancelled:     tm.CanBeCancelled,
			MoveProposedBy:     tm.MoveProposedBy,
			MoveRequiresAccept: tm.CanBeCancelled,
			Notes:              tm.Notes,
			ProposedTime:       tm.ProposedTime,
			Time:               tm.Time,
			User:               tm.User,
			UserUuid:           tm.UserUUID,
			Uuid:               tm.UUID,
		}

		trainings = append(trainings, t)
	}

	return trainings
}

func newDomainUserFromAuthUser(ctx context.Context) (training.User, error) {
	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		return training.User{}, err
	}

	userType, err := training.NewUserTypeFromString(user.Role)
	if err != nil {
		return training.User{}, err
	}

	return training.NewUser(user.UUID, userType)
}
