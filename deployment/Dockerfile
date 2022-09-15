FROM scratch

COPY /metrics-multiplexer-sidecar /metrics-multiplexer-sidecar

# Max user
USER 255999


EXPOSE 13434
ENTRYPOINT ["/metrics-multiplexer-sidecar"]
