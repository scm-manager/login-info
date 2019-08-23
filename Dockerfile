FROM scratch
COPY target/login-info /login-info
USER 10000
CMD ["/login-info"]
