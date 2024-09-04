FROM golang:1.23

WORKDIR /

COPY ./ ./

RUN go mod download
RUN go build -o forum ./cmd/forum/main.go

EXPOSE 8000

CMD ["./forum"]