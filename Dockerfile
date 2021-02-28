FROM golang:1.15

WORKDIR /app

COPY go.mod go.sum entrypoint.sh ./
RUN go mod download
COPY app.go ./

CMD ["./entrypoint.sh"]

EXPOSE 8080
