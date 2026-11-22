# Build stage
FROM golang:1.22 AS base

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o main .

# Runtime stage
FROM gcr.io/distroless/base-debian13:nonroot AS final

WORKDIR /app

# Same for one line of code
#COPY --from=base /app/main /app/static /app/templates ./

COPY --from=base /app/main .
COPY --from=base /app/static ./static
COPY --from=base /app/templates ./templates

EXPOSE 8080

CMD ["./main"]