FROM scratch
ARG binary
ARG version
ENV version=$version
ADD bin/$binary /app

expose 8090
ENTRYPOINT ["/app"]