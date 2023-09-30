# The build stage
FROM golang:1.21.1-alpine3.18 AS build

WORKDIR /app

# copy the src code into the /app folder
COPY . .

# build the image from this .go file and name the exe output as paywise.exe which will be stored within the /app folder
RUN go build -o paywise ./cmd/main.go

# The Run stage
FROM alpine:3.18

WORKDIR /app

# copy from the build stage the paywise exe file which has a path = /app/paywise into my current workdir which is the /app folder which is defined in the new stage (brand new /app folder)
COPY --from=build /app/paywise .

# Expose this port within my container (just for documentation)
EXPOSE 8000

# Default command to be run when the container starts
CMD ["/app/paywise"]
