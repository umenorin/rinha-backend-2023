FROM golang:1.24

ENV POSTGRES_URL=postgres://postgres:mysecretpassword@localhost:5432/rinha?sslmode=disable
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app  ./cmd/client/ 


CMD ["app"]

