FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN GOFLAGS=-mod=mod go mod tidy && CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app-service .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /app-service /app-service
EXPOSE 3000
ENTRYPOINT ["/app-service"]
