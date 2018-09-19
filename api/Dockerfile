FROM golang:1.11.0

WORKDIR /go/src/github.com/knixdev/domain-info/api
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["api"]