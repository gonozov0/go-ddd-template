FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY ./bin/your-service-name /your-service-name

CMD ["/your-service-name"]