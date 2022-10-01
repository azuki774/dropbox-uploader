FROM alpine:3.7
ADD build/bin/dropbox-uploader /usr/local/bin/dropbox-uploader
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

ENTRYPOINT ["/usr/local/bin/dropbox-uploader", "start"]

