package httperr_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	commonErrors "github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/errors"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server/httperr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// customSlugError is an error with its own type, like the ones defined in the domain layers.
type customSlugError struct{}

func (customSlugError) Error() string { return "training time is too close" }

func (customSlugError) Slug() string { return "training-time-too-close" }

func (customSlugError) ErrorType() commonErrors.ErrorType {
	return commonErrors.ErrorTypeIncorrectInput
}

func TestRespondWithSlugError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		Error           error
		ExpectedStatus  int
		ExpectedSlug    string
		ExpectedMessage string
	}{
		{
			Name:            "incorrect_input",
			Error:           commonErrors.NewIncorrectInputError("Note too long", "note-too-long"),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedSlug:    "note-too-long",
			ExpectedMessage: "Note too long",
		},
		{
			// domain errors are usually returned wrapped, so unwrapping needs to work
			Name:            "wrapped_incorrect_input",
			Error:           errors.WithStack(commonErrors.NewIncorrectInputError("Note too long", "note-too-long")),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedSlug:    "note-too-long",
			ExpectedMessage: "Note too long",
		},
		{
			Name:            "authorization",
			Error:           commonErrors.NewAuthorizationError("no user in context", "no-user-found"),
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedSlug:    "no-user-found",
			ExpectedMessage: "no user in context",
		},
		{
			Name:            "custom_error_type",
			Error:           errors.WithStack(customSlugError{}),
			ExpectedStatus:  http.StatusBadRequest,
			ExpectedSlug:    "training-time-too-close",
			ExpectedMessage: "training time is too close",
		},
		{
			// errors that were not written for the user must not leak their message
			Name:            "unknown_error",
			Error:           errors.New("some internal failure"),
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedSlug:    "internal-server-error",
			ExpectedMessage: "Internal server error",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			resp := doRequest(t, func(w http.ResponseWriter, r *http.Request) {
				httperr.RespondWithSlugError(c.Error, w, r)
			})
			defer func() {
				require.NoError(t, resp.Body.Close())
			}()

			assert.Equal(t, c.ExpectedStatus, resp.StatusCode)

			// without it clients don't parse the body, so the error can't be displayed
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var errorResponse httperr.ErrorResponse
			require.NoError(t, json.Unmarshal(body, &errorResponse))

			assert.Equal(t, c.ExpectedSlug, errorResponse.Slug)
			assert.Equal(t, c.ExpectedMessage, errorResponse.Message)
		})
	}
}

// doRequest calls the handler over a real connection: an httptest.ResponseRecorder keeps its
// headers writable after WriteHeader, so it wouldn't catch headers lost after the status is sent.
func doRequest(t *testing.T, handler http.HandlerFunc) *http.Response {
	t.Helper()

	logger := logrus.New()
	logger.Out = io.Discard

	server := httptest.NewServer(logs.NewStructuredLogger(logger)(handler))
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL)
	require.NoError(t, err)

	return resp
}
