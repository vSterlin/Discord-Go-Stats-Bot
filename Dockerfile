FROM golang:1.14

WORKDIR vs/go/src/statsbot
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["statsbot"]