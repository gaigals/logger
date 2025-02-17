FROM golang:1.24

COPY *.go /app/
COPY go.mod /app/

COPY example/go.mod /app/example/go.mod
COPY example/main.go /app/example/main.go


WORKDIR /app/example

CMD ["go", "run", "./"]
