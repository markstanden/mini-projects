FROM golang:1.16.3 AS buildstep

# install git to pull down repo
RUN apt-get install git
#RUN psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" --port "9000" -c "test"

# clone the repo
WORKDIR /app
RUN git clone https://github.com/markstanden/authentication.git /app

# download the dependancies
RUN go mod download

# Run our tests, verbose output
RUN CGO_ENABLED=0 go test ./...

# Build the executable
RUN go build -o ./build .


#
#   Now create the minimal run container
#
FROM alpine:latest
#need to work on the certs later
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
#COPY ./mycert.crt /usr/local/share/ca-certificates/mycert.crt
#RUN update-ca-certificates

COPY --from=buildstep /app/build/ /app

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/authentication"]