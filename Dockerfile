FROM golang:1.12

WORKDIR /code
ENTRYPOINT ["go", "run", "main.go"]