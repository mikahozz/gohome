# base image
FROM golang:1.19.2-alpine3.16 as base
WORKDIR /builder
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /builder/main /builder/main.go

# runner image
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /app
COPY --from=base /builder/main main

EXPOSE 6001
CMD ["/app/main"]