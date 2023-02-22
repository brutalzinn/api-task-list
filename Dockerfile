FROM golang:1.20.1

# Set the current working directory inside the container
WORKDIR /app

ENV TZ=America/Sao_Paulo

RUN go install github.com/cosmtrek/air@latest

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the workspace
# COPY . .

# Build the Go app
# RUN go build -o main .

# Command to run the executable
CMD ["air", "-c", ".air.docker.toml"]
