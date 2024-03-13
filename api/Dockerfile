FROM golang:alpine AS build

WORKDIR /app

ENV GO111MODULE=on
COPY api/go.mod api/go.sum ./
RUN go mod download

COPY /api/ ./
RUN go build -o main cmd/main.go

EXPOSE 8080

FROM scratch

COPY --from=build app/main /bin/main

ENTRYPOINT ["/bin/main"]