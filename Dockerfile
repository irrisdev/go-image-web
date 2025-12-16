FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o go-image-web main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /app/go-image-web .
COPY --from=build /app/public ./public
COPY --from=build /app/internal/db/migrations ./internal/db/migrations

EXPOSE 9991

ENTRYPOINT ["./go-image-web"]