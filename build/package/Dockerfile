FROM golang:1.18 as gobuild
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 go build -o /migrate $(head -1 go.mod | awk '{print $2}')/cmd/migrate

FROM scratch
ENTRYPOINT ["/migrate"]
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=gobuild /migrate /migrate
