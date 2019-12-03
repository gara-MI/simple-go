############################
# STEP 1 use builder to generate go binary
############################
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

# Create appuser.
RUN adduser -D -g '' axway

WORKDIR $GOPATH/src/github/gara-MI/simple-go/

COPY . .

# Fetch dependencies. using go get.
RUN go get -d -v

# Build the binary. CGO_ENABLED=0 compulosry for scratch image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/main

############################
# STEP 2 build a small image
############################
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable.
COPY --from=builder /go/bin/main /go/bin/main

# Use an unprivileged user.
USER axway

ENTRYPOINT ["/go/bin/main"]
