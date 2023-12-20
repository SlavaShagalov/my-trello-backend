FROM golang:1.20-alpine3.18 AS install
RUN apk add --no-cache git make bash
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go install github.com/mailru/easyjson/...@latest
WORKDIR /src
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download

FROM install AS copy
WORKDIR /src
ENV ALLURE_OUTPUT_PATH=/src
COPY internal ./internal
COPY tests ./tests
COPY scripts ./scripts
COPY make ./make
COPY Makefile ./
RUN make mocks
RUN make easyjson
CMD ["tail", "-f", "/dev/null"]
