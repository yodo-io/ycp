## Builder
FROM golang:1.11-alpine3.8 AS builder

ENV GO111MODULE=on

RUN apk --no-cache add git build-base
WORKDIR /src/
ADD . ./

RUN go test ./...
RUN GOOS=linux go build -a -installsuffix cgo -o /bin/ycp .


## Runtime image
FROM alpine:3.8
EXPOSE 9000
COPY --from=builder /bin/ycp /bin
CMD ["/bin/ycp"]
