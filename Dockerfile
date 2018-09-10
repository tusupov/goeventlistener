FROM golang:1.11

# copy project
WORKDIR /go/src/github.com/tusupov/goeventlistener
COPY . ./

# run test
RUN go test -v ./...
RUN go test --bench=. -v ./...
RUN rm -rf *_test.go

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o eventlistener .

FROM alpine:latest

# add certificates for https connections
RUN apk --no-cache add ca-certificates

# copy
WORKDIR /app/
COPY --from=0 /go/src/github.com/tusupov/goeventlistener/eventlistener .

CMD ./eventlistener -p $PORT
