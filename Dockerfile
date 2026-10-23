FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache make
RUN make install
CMD ["cheatsh"]
