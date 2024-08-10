# base image
FROM golang:alpine3.20 as base
WORKDIR /builder
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /builder/main /builder/main.go

# runner image
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /app
COPY --from=base /builder/main /builder/.env ./

EXPOSE 6001
CMD ["/app/main"]