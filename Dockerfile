FROM golang:alpine as builder

WORKDIR /app
COPY . /app
RUN go build -o /app/bin/worker ./worker
RUN go build -o /app/bin/test-me ./test-me 

FROM golang:alpine

COPY --from=builder /app/bin/worker /app/bin/worker
COPY --from=builder /app/bin/test-me /app/bin/test-me
