FROM golang:1.16-buster
WORKDIR /go/src/app
RUN git clone https://github.com/jasonbronson/router-device-tracker .

RUN apt-get install gcc make
RUN go get -d -v ./...
RUN go install -v ./...
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./routerapp /go/src/app/*.go

CMD ["routerapp"]