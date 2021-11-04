FROM golang:1.17.0

WORKDIR /arya-core

ADD . /arya-core

#RUN go build -o main .

CMD ["/bin/bash"]
