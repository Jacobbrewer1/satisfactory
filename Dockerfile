FROM golang:1.21.1

LABEL org.opencontainers.image.source='https://github.com/Jacobbrewer1/satisfacotry'
LABEL org.opencontainers.image.description="A package that is used around the Satisfactory game"
LABEL org.opencontainers.image.licenses='GNU General Public License v3.0'

WORKDIR /cmd

# Copy the binary from the build
COPY ./bin/app /cmd/app

RUN ["chmod", "+x", "./app"]

ENTRYPOINT ["/cmd/app"]
