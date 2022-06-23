FROM golang:1.17

# ENV GO111MODULE=on
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src

RUN go get github.com/joho/godotenv
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/go-redis/redis/v8 
# RUN go mod download
RUN cd /go/src && go build .

CMD ["go", "run", "/go/src/main.go"]
