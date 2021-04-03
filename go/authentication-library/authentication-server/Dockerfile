# FROM postgres AS db
# RUN apt-get update && apt-get install procps -y

# #RUN createdb test

# #RUN apt-get update && apt-get install -q -y postgresql-9.6 postgresql-client-9.6 postgresql-contrib-9.6 postgresql-client-common postgresql-common
# RUN echo postgres:postgres | chpasswd
# RUN pg_createcluster 13 main --start
# RUN pg_ctlcluster 13 main start
# #RUN /etc/init.d/postgresql start
 
# RUN ps -ef | grep postgres
# RUN find /var/run/postgresql/ -name .s.PGSQL.5432
# EXPOSE 5432

# RUN pg_lsclusters

# RUN cat /var/log/postgresql/postgresql-13-main.log

# #RUN psql -c "ALTER USER postgres PASSWORD 'postgres';" postgres

#RUN psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" --port "5432" -c "test"

FROM golang:1.16.3 AS builder

# install git to pull down repo
RUN apt-get install git -y

ENV PGTESTHOST=localhost

# clone the repo
WORKDIR /app
RUN git clone https://github.com/markstanden/authentication.git /app

# download the dependancies
RUN go mod download

ENV CGO_ENABLED=0
# Run our tests, verbose output
CMD go test ./...

# Build the executable
#RUN go build -o ./build/ ./cmd/authentication


#
#   Now create the minimal run container
#
#FROM alpine:latest
#need to work on the certs later
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
#COPY ./mycert.crt /usr/local/share/ca-certificates/mycert.crt
#RUN update-ca-certificates

#COPY --from=buildstep /app/build/ /app

# This container exposes port 8080 to the outside world
#EXPOSE 8080
#ENV PORT=3000
# Run the binary program produced by `go install`
#CMD ["./build/authentication"]