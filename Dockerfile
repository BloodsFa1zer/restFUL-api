FROM golang:1.20-alpine as dev-env

# Copy application data into image
COPY . /Users/mishashevnuk/GolandProjects/app3.1
WORKDIR /Users/mishashevnuk/GolandProjects/app3.1


COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /server


##
## Deploy
##
FROM alpine:latest
RUN mkdir /data

COPY --from=dev-env /server ./
CMD ["/server"]
