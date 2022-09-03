FROM golang:alpine AS builder
ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR "/srv"
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY ./whistle-pig /srv/whistle-pig
USER appuser:appuser
EXPOSE 8088
ENTRYPOINT ["/srv/whistle-pig"]
