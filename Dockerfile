# First stage: Build the Go binary
FROM golang:1.23.6 AS builder

WORKDIR /exporter

ADD ./ /exporter


RUN make build/exporter

# Second stage: Create a minimal runtime image
FROM debian:bookworm-slim

WORKDIR /

# Copy built binaries from the builder stage
COPY --from=builder /exporter/bin/polkadot-exporter-local-compatible /exporter/polkadot-exporter

# Expose the necessary ports
EXPOSE 9100

# Set the entrypoint for the container
ENTRYPOINT ["/exporter/polkadot-exporter"]
