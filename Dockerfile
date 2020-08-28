# Compile stage
FROM golang:1.13.8 AS build-env
# ADD . /dockerdev
COPY . /usr/local/go/src/librarian/
WORKDIR /usr/local/go/src/librarian/
CMD ["go", "get", "zombiezen.com/go/capnproto2"]
# CMD ["cp", "test.so", "/test.so"]
RUN go build -o /app

# Final stage
FROM debian:buster
EXPOSE 8000
WORKDIR /
COPY --from=build-env /app /
COPY --from=build-env /usr/local/go/src/librarian/gosharedlib.so /
CMD ["/app"]