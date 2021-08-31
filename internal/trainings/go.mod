module github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings

go 1.14

require (
	cloud.google.com/go v0.38.0
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common v0.0.0-00010101000000-000000000000
	github.com/deepmap/oapi-codegen v1.4.1
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.5.0
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	google.golang.org/api v0.21.0
	google.golang.org/grpc v1.40.0
)

replace github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common => ../common/
