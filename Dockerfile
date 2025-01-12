FROM golang:1.23-alpine AS build

RUN apk add --no-cache curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod

WORKDIR /app
COPY --from=build /app/main /app/main
COPY templates/ /app/templates/
COPY assets/ /app/assets/

EXPOSE ${PORT}
CMD ["./main"]
