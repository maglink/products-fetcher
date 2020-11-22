FROM golang:1.15.2
COPY ./ /app
WORKDIR /app
RUN go build -o ./bin/products-fetcher ./cmd/products-fetcher
CMD [ "/app/bin/products-fetcher" ]