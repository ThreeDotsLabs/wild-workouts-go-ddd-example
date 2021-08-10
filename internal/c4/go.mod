module github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/c4

go 1.16

require (
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer v0.0.0-00010101000000-000000000000
	github.com/krzysztofreczek/go-structurizr v0.1.2
)

replace (
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common => ../common/
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer => ../trainer/
)
