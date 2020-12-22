FROM golang:1.15 as backend-build

WORKDIR /membership
COPY . .

RUN go mod vendor
RUN go build -o server

ENTRYPOINT [ "./server" ]
