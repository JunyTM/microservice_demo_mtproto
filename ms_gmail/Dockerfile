FROM golang:1.22

WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y \
    git \
    && rm -rf /var/lib/apt/lists/*

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# RUN chmod +rw excel

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o out/ms_gmail ./

EXPOSE 8080

CMD ["out/ms_gmail"]