FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-image-web main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /app/go-image-web .
COPY ./public /app/public


EXPOSE 9991

ENTRYPOINT ["./go-image-web"]