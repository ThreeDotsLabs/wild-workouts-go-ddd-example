package httperr

import (
	stderrors "errors"
	"net/http"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/errors"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/go-chi/render"
)

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	// internal errors are not written for the user, so their message is never passed to the client
	httpRespondWithError(err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

func Unauthorised(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, userMessage(err, "Unauthorised"), http.StatusUnauthorized)
}

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, userMessage(err, "Bad request"), http.StatusBadRequest)
}

func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	var slugError errors.SlugErrorer
	if !stderrors.As(err, &slugError) {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorised(slugError.Slug(), err, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), err, w, r)
	default:
		InternalError(slugError.Slug(), err, w, r)
	}
}

// userMessage returns the message of errors that were written to be shown to the user
// (the ones implementing errors.SlugErrorer), and a generic fallback for all the others.
func userMessage(err error, fallback string) string {
	var slugError errors.SlugErrorer
	if stderrors.As(err, &slugError) {
		return slugError.Error()
	}

	return fallback
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, message string, status int) {
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(message)
	resp := ErrorResponse{slug, message, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Message    string `json:"message"`
	httpStatus int
}

// Render sets the status via the request context instead of calling w.WriteHeader directly:
// the response headers are frozen on the first WriteHeader, so writing it here would drop the
// "Content-Type: application/json" that render sets afterwards, and clients wouldn't parse the body.
func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.httpStatus)
	return nil
}
