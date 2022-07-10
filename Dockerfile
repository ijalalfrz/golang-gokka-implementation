# Image Builder
FROM telkomindonesia/alpine:go-1.15 AS go-builder

LABEL maintainer="ijal.alfarizi@gmail.com"

# Set Working Directory
WORKDIR /usr/src/app

# Copy Source Code
COPY . ./

# Dependencies installation and binary file builder
RUN make install \
  && make build


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base

# Set Working Directory
WORKDIR /usr/src/app

# Copy Anything The Application Needs
COPY --from=go-builder /tmp/app ./


# Remove the command below if the service doesn't need JWT Signer / Parser
# Check the Makefile
COPY --from=go-builder /tmp/secret secret

# Expose Application Port
EXPOSE 9000

# Run The Application
CMD ["./app"]