# The base stage
FROM golang:1.21.1-alpine3.18

WORKDIR /app

# copy the src code into the /app folder
COPY . .

# build the image from this .go file and name the exe output as paywise.exe which will be stored within the /app folder
RUN go build -o paywise ./cmd/main.go

# Expose this port within my container (just for documentation)
EXPOSE 8000

# Default command to be run when the container starts
CMD ["/app/paywise"]
