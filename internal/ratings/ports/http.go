package ports

import (
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/ratings/app"
	"github.com/go-chi/render"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) GetRatings(w http.ResponseWriter, r *http.Request) {
	// TODO: add real implementation, you can check internal/trainer/ports/http.go for the inspiration
	render.Respond(w, r, Ratings{Ratings: []Rating{
		{
			CanRate: false,
			Date:    timePtr(time.Date(2021, 10, 3, 11, 1, 0, 0, time.Local)),
			Comment: strPtr("best trainer in the world! 💕😍"),
			Rated:   true,
			Rating:  intPtr(5),
		},
		{
			CanRate: false,
			Date:    timePtr(time.Date(2021, 10, 5, 22, 7, 0, 0, time.Local)),
			Comment: strPtr("watching Netflix is better"),
			Rated:   true,
			Rating:  intPtr(1),
		},
		{
			CanRate: false,
			Date:    timePtr(time.Date(2021, 11, 6, 19, 15, 0, 0, time.Local)),
			Comment: strPtr("it was pretty cool cool, but you can swear a bit less during workout 😅"),
			Rated:   true,
			Rating:  intPtr(4),
		},
		{
			CanRate: false,
			Date:    timePtr(time.Date(2021, 11, 8, 12, 5, 0, 0, time.Local)),
			Rated:   true,
			Rating:  intPtr(5),
		},
	}})
}

func (h HttpServer) GetTrainingRating(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	// TODO: add real implementation, you can check internal/trainer/ports/http.go for the inspiration
	render.Respond(w, r, Rating{
		CanRate: true,
		Rated:   false,
	})
}

func (h HttpServer) PostTrainingRating(w http.ResponseWriter, r *http.Request, trainingUUID string) {
	// TODO: add real implementation, you can check internal/trainer/ports/http.go for the inspiration
	w.WriteHeader(http.StatusNoContent)
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
