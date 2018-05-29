FROM golang:1.10 as gobuild
WORKDIR /app
RUN go get -u golang.org/x/vgo
COPY . ./
RUN CGO_ENABLED=0 vgo build -o /migrate $(head -1 go.mod | awk '{print $2}')/cmd/migrate

FROM scratch
CMD ["/migrate"]
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=gobuild /migrate /migrate
