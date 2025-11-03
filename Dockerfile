############# builder
FROM golang:1.25.3 AS builder

WORKDIR /go/src/github.com/MartinWeindel/gardener-extension-dummy-service

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

ARG EFFECTIVE_VERSION
RUN make install EFFECTIVE_VERSION=$EFFECTIVE_VERSION

############# gardener-extension-dummy-service
FROM gcr.io/distroless/static-debian12:nonroot AS gardener-extension-dummy-service
WORKDIR /

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-dummy-service /gardener-extension-dummy-service
ENTRYPOINT ["/gardener-extension-dummy-service"]
