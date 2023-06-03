FROM alpine:3.18
LABEL org.opencontainers.image.source=https://github.com/babadro/forecaster

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup --system app && adduser --system app --ingroup app

COPY release/app /usr/bin/app

RUN chmod +x /usr/bin/app

USER app

CMD ["/usr/bin/app"]