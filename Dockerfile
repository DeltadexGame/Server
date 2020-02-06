FROM golang:1.13

WORKDIR /deltadex

ENV GO111MODULES=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build
ENTRYPOINT ["./deltadex"]