FROM alpine:3.18
LABEL org.opencontainers.image.source=https://github.com/babadro/forecaster

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup --system app && adduser --system app --ingroup app

COPY release/app /usr/bin/service

RUN chmod +x /usr/bin/service

USER app

CMD ["service"]