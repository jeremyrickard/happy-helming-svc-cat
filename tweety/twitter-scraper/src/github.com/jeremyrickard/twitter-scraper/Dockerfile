from golang:1.9

WORKDIR /go/src/tweets
COPY . .

RUN go install -v ./...

CMD ["../../bin/tweets"]
