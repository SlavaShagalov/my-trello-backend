FROM golang:1.20-alpine3.18 AS install
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS build
WORKDIR /src
COPY cmd ./cmd
COPY internal ./internal
COPY docs ./docs
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /bin/api cmd/api/main.go

FROM ubuntu AS api
WORKDIR /
COPY --from=build /bin/api /bin/api
CMD ["/bin/api"]
