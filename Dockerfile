FROM alpine:3.7
ARG BINARY
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY ${BINARY} /app/license-manager
COPY user_licenses.db /app
WORKDIR /app
ENTRYPOINT ["./license-manager"]
