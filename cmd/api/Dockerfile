FROM golang:1.19.6-alpine3.17 AS install
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS build
WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/api cmd/api/main.go

FROM scratch AS api
WORKDIR /
COPY --from=build /bin/api /bin/api
CMD ["/bin/api"]
