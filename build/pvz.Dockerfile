FROM golang:1.24.0-alpine AS builder

COPY go.mod go.sum /github.com/marrgancovka/pvzService/
WORKDIR /github.com/marrgancovka/pvzService/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/pvz/main.go

FROM scratch AS runner

WORKDIR /docker-cian/

COPY --from=builder /github.com/marrgancovka/pvzService/.bin .
COPY --from=builder /github.com/marrgancovka/pvzService/config config/

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

ENTRYPOINT ["./.bin"]