FROM golang:latest
WORKDIR /app
COPY ./src .
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -command="go run main.go" -exclude-dir=.git -polling
EXPOSE 80