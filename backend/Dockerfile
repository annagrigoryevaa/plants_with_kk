FROM golang:1.25 AS build

WORKDIR /src/backend

COPY backend/go.mod ./
RUN go mod download

COPY backend/. ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/plant-keeper .

FROM debian:bookworm-slim

WORKDIR /app/backend

COPY --from=build /out/plant-keeper /usr/local/bin/plant-keeper
COPY frontend /app/frontend

RUN mkdir -p /data

ENV PORT=8080
ENV DATA_DIR=/data

EXPOSE 8080

CMD ["plant-keeper"]
