# build stage
FROM golang:1.22-bookworm AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o ./app ./cmd/api/main.go

# deploy stage
FROM gcr.io/distroless/base-debian12 AS deploy

WORKDIR /app

COPY --from=build ./app/app ./app

EXPOSE 8080

# needed for distroless base
USER nonroot:nonroot

ENTRYPOINT [ "./app" ]
