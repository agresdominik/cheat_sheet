FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache make
RUN make install
CMD ["cheatsh"]
