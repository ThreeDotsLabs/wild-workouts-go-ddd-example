module github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/users

go 1.14

require (
	cloud.google.com/go/firestore v1.2.0 // indirect
	github.com/deepmap/oapi-codegen v1.3.6 // indirect
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/render v1.0.1 // indirect
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common v0.0.0-00010101000000-000000000000
)

replace github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common => ../common/
