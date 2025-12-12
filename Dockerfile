FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-image-web main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=build /app/go-image-web .

EXPOSE 9991

ENTRYPOINT ["./go-image-web"]