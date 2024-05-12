FROM golang:1.22 as builder

RUN mkdir /app
ENV GO111MODULE=on

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server /app/cmd/goshell/.

FROM scratch
COPY --from=builder /app/server /
COPY --from=builder /app/config/local.yaml /
EXPOSE 8080
CMD ["/server"]