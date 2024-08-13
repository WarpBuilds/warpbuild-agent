# Dockerfile
FROM alpine:latest

# Example operation to trigger a cache hit/miss
RUN apk update && apk add --no-cache curl

CMD ["echo", "Hello World"]