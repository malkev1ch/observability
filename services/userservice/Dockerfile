# Docker images to build application.
ARG BUILDER_IMAGE=golang:1.21.3

FROM ${BUILDER_IMAGE}
WORKDIR /build
COPY . .

RUN go mod download
RUN cd userservice && go build -o /bin/userservice

# The command serves as a record, specifying the ports that the container is intended to operate on.
EXPOSE 8080

ENTRYPOINT ["/bin/userservice"]
CMD []
