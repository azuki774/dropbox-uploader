FROM alpine:3.7
ADD build/bin/dropbox-uploader /usr/local/bin/dropbox-uploader
# Copy source directory set /data/
ENTRYPOINT ["/usr/local/bin/dropbox-uploader", "upload"]
CMD ["-t", "<token>", "-s", "<src-dir> ", "-d", "<dst-dir>"]
