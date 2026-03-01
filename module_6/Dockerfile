FROM golang:1.25-alpine3.23 AS build

WORKDIR /work

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o server ./cmd/server

FROM alpine:3.23

WORKDIR /work

COPY --from=build /work/server .

EXPOSE 8080

CMD ["./server"]
