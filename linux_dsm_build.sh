#docker build ./docker/nas --tag sound

docker run --rm -v "$PWD":/go/src/github.com/hednowley/sound -w /go/src/github.com/hednowley/sound golang:latest go get -d -v && GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -v
