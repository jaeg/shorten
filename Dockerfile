FROM scratch
ARG binary
ARG version
ENV version=$version
ADD bin/$binary /app
add config config

expose 8090
ENTRYPOINT ["/app"]