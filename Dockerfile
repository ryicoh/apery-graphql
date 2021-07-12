FROM golang:1.16-buster AS builder

WORKDIR /src

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go build -o apery-graphql ./cmd/main.go


FROM ryicoh/apery

COPY --from=builder /src/apery-graphql /app/qpery-graphql

CMD "/app/qpery-graphql" "--binary" "./apery"
