FROM golang:alpine AS build

WORKDIR /app

COPY /api/ ./
RUN go build -o main cmd/main.go

EXPOSE 8080

FROM scratch

COPY --from=build app/main /bin/main

ENTRYPOINT ["/bin/main"]