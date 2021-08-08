# Production image based on alpine.
FROM alpine
LABEL maintainer="Lukas Malkmus <mail@lukasmalkmus.com>"

# Upgrade packages and install ca-certificates.
RUN apk update --no-cache
RUN apk upgrade --no-cache
RUN apk add --no-cache ca-certificates

# Copy binary into image.
COPY horcrux /usr/bin/horcrux

# Use the project name as working directory.
WORKDIR /horcrux

# Set the binary as entrypoint.
ENTRYPOINT [ "/usr/bin/horcrux" ]