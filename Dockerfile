FROM golang:latest
WORKDIR /app

COPY . .
RUN go build -o updowninserve .

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/updowninserve /bin/
COPY index.html .
COPY htmx.min.js .

CMD ["/bin/updowninserve"]

