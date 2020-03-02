FROM golang:latest 
ENV PORT=3000
RUN mkdir /go/src/knowledgebase 
ADD . /go/src/knowledgebase 
WORKDIR /go/src/knowledgebase
RUN go get
RUN go install
RUN go build -o main .

EXPOSE 3000
CMD ["/go/src/knowledgebase/main"]