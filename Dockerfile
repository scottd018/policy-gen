## NOTE: This image uses goreleaser to build an image
# if building manually please run go build ./internal/cmd/policygen first and then build

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:latest as production

COPY policy-gen /usr/bin/policy-gen
RUN chmod +x /usr/bin/policy-gen

WORKDIR /workdir

ENTRYPOINT ["/usr/bin/policy-gen"]
CMD ["--help"]