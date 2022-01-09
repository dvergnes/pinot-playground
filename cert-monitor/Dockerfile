FROM gcr.io/distroless/static-debian11:non

COPY --chown=nonroot:nonroot /bin/cert-monitor_linux_amd64 /bin/cert-monitor

ENTRYPOINT ["/bin/cert-monitor", "--config=/config/config.yml"]
