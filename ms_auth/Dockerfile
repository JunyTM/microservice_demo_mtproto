FROM golang:1.22

WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y \
    git \
    && rm -rf /var/lib/apt/lists/*

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /out/main ./

# wait-for mysql to be installed before running the app in the background
# NOTE: this is will avoid a error in the background: "Connection refused"
RUN wget -O wait-for https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
RUN chmod +x wait-for
RUN mv wait-for /usr/local/bin/

ENTRYPOINT ["wait-for", "authen_db:3306", "--", "/out/main", "-db", "true"]