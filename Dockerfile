FROM golang:1.16-alpine AS builder
RUN apk --no-cache add ca-certificates git
WORKDIR /app

ENV GIN_MODE=release

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN go build -o ./api ./main.go

# Create final image
FROM alpine AS runner
WORKDIR /

ENV GIN_MODE=release

# Copy config file & complied file
COPY --from=builder /app/.env .
COPY --from=builder /app/api .

EXPOSE 31337
CMD ["./api"]
