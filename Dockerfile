FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/vngitPub

FROM alpine:latest
WORKDIR /apprun
COPY --from=builder /app/vngitPub /apprun/vngitPub
COPY --from=builder /app/VERSION /apprun/VERSION
EXPOSE 8000
CMD ["/apprun/vngitPub"]