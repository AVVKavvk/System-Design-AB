FROM golang:1.26

WORKDIR /usr/src/pg-bouncer

COPY . .

RUN go mod download

RUN go build -o main .

CMD [ "./main" ]