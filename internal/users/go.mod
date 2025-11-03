module github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/users

go 1.18

require (
	cloud.google.com/go/firestore v1.5.0
	firebase.google.com/go/v4 v4.7.1
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.5.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/api v0.40.0
	google.golang.org/grpc v1.40.0
)

require (
	cloud.google.com/go v0.75.0 // indirect
	cloud.google.com/go/storage v1.10.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/cors v1.0.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.opencensus.io v0.22.5 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210218202405-ba52d332ba99 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210222152913-aa3ee6e6a81c // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

replace github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common => ../common/
