FROM golang:alpine

RUN apk add --no-cache curl bash ca-certificates \
 && update-ca-certificates

Add . $GOPATH/src/poller
Run go install poller
Run chmod +x $GOPATH/src/poller/run.sh
Run cp $GOPATH/src/poller/run.sh /go/bin/run.sh

CMD ["run.sh"]
