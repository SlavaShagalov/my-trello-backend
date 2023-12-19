FROM golang:1.20-alpine3.18 AS install
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS copy
WORKDIR /src
ENV ALLURE_OUTPUT_PATH=/src
COPY internal ./internal
COPY tests ./tests
COPY make ./make
COPY Makefile ./
CMD ["tail", "-f", "/dev/null"]
