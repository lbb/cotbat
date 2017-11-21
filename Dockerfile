FROM alpine:latest
ADD main /
WORKDIR /
CMD ./main
