module github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/users

go 1.17

require (
	cloud.google.com/go/firestore v1.2.0
	firebase.google.com/go v3.12.0+incompatible
	github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.5
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.5.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	google.golang.org/api v0.21.0
	google.golang.org/grpc v1.40.0
)

require (
	cloud.google.com/go v0.55.0 // indirect
	cloud.google.com/go/storage v1.6.0 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-chi/cors v1.0.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.2.0 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20211031064116-611d5d643895 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.0.0-20200325010219-a49f79bcc224 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
)

replace github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common => ../common/
