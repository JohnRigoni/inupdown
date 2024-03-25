FROM golang:latest
WORKDIR /app

COPY . .
RUN go build -o inupdown .

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/inupdown /bin/
COPY index.html .
COPY htmx.min.js .

CMD ["/bin/inupdown"]

