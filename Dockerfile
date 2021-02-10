FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go build -v -o bin/golang-game-server ./cmd/golang-game-server/

CMD [ "/go/src/app/bin/golang-game-server" ]