
FROM golang:alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

RUN apk add --no-cache ca-certificates git

WORKDIR /src

# COPY ./go.mod ./go.sum ./
# RUN go mod download

COPY ./*.png ./*.xml ./*.json ./*.go ./*.ico ./
RUN go mod init github.com/Zate/HashHound
RUN go mod tidy

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

WORKDIR /

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /app

COPY --from=builder /src/*.png /

COPY --from=builder /src/browserconfig.xml .

COPY --from=builder /src/favicon.ico .

COPY --from=builder /src/manifest.json .

EXPOSE 3003

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]