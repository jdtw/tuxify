FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /app ./...

FROM gcr.io/distroless/static:latest
WORKDIR /app
COPY --from=builder /app/tuxify-server /app/tuxify-server

EXPOSE 8080
CMD ["/app/tuxify-server"]