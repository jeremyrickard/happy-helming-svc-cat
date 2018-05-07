from golang:1.9

WORKDIR /go/src/tweety
COPY . .

RUN go install -v ./...

CMD ["../../bin/tweety"]
