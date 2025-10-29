FROM golang:1.25

RUN go install github.com/cespare/reflex@v0.3.1
COPY reflex.conf /

COPY start.sh /
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
