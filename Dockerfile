FROM golang:1.12.7 as builder
WORKDIR /go/src/bitbucket.org/scm-manager/login-info
COPY . .
RUN GO111MODULE=on go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o login-info *.go

FROM scratch
COPY --from=builder /go/src/bitbucket.org/scm-manager/login-info/login-info /login-info
USER 10000
CMD ["/login-info"]
