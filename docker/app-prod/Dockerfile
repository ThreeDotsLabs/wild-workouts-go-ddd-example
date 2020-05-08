FROM golang:1.14

ARG SERVICE

RUN mkdir /internal
COPY . /internal
WORKDIR /internal/$SERVICE
RUN go build -o /app .

CMD ["/app"]
