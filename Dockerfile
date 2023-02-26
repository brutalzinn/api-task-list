FROM golang:latest

ENV TZ=America/Sao_Paulo
WORKDIR /app
ADD . /app

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
COPY go.mod go.sum ./
RUN go mod download

ENTRYPOINT CompileDaemon --build="go build main.go" -exclude-dir=.git -polling