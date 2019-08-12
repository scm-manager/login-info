FROM golang:1.12.7 as builder
WORKDIR /go/src/bitbucket.org/scm-manager/login-info
COPY . .
# RUN GO111MODULES=on go get -u
RUN GO111MODULE=on go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o login-info *.go

FROM scratch
COPY --from=builder /go/src/bitbucket.org/scm-manager/login-info/login-info /login-info
CMD ["/login-info"]
