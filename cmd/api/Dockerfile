FROM golang:1.21-alpine3.18 AS install
RUN apk add --no-cache git make bash
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go install github.com/mailru/easyjson/...@latest
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS build
WORKDIR /src
COPY cmd ./cmd
COPY internal ./internal
COPY docs ./docs
COPY scripts ./scripts
COPY make ./make
COPY Makefile ./
RUN make mocks
RUN make easyjson
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /bin/api cmd/api/main.go

FROM ubuntu AS api
WORKDIR /
COPY --from=build /bin/api /bin/api
CMD ["/bin/api"]
