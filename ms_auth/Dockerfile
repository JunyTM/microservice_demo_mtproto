FROM golang:1.22

WORKDIR /usr/src/app



# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /out/main ./

# Add an entry to /etc/hosts
COPY modify_hosts.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/modify_hosts.sh

ENTRYPOINT ["/out/main"]
