#############################################################################
####### Stage 1 - Build

FROM golang:1.13-rc as builder

#ARG APP_NAME=vegeta-kubernetes
ARG GOPATH=/go

LABEL maintainer="Tamir Aviv"

WORKDIR /go/src/vegeta-kubernetes

COPY . .

RUN go get -d ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/vegeta-kubernetes /go/src/vegeta-kubernetes/cmd/vegeta-kubernetes/main.go

#############################################################################
####### Stage 2 - Run from scratch

FROM scratch

COPY --from=builder /go/bin/vegeta-kubernetes .

EXPOSE 8080

CMD ["./vegeta-kubernetes"] 