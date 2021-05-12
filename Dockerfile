FROM golang:latest
COPY . /app
WORKDIR /app
RUN GOOS=linux go build -a -ldflags="-linkmode external -extldflags -static" -o main && chmod +x main

FROM gcr.io/distroless/static
COPY --from=0 /app/main /main
CMD ["/main"]

