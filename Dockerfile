FROM alpine:3.24

RUN addgroup -S kurz && adduser -S kurz -G kurz
USER kurz

COPY bin/kurz /srv/kurz

EXPOSE 8080

ENTRYPOINT ["/srv/kurz"]
