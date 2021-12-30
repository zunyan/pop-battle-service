FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 5000

CMD ["go", "run", "main.go"]