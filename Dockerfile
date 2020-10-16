FROM golang:latest
ADD . /go/src/github.com/madflojo/efs2
WORKDIR /go/src/github.com/madflojo/efs2
RUN go install -v
ENTRYPOINT ["efs2"]
