FROM golang:1.22.1-alpine3.19 as builder


WORKDIR /app


COPY go.mod .
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o converter ./cmd


FROM alpine:latest  


RUN apk --no-cache add ca-certificates pandoc

WORKDIR /root/


COPY --from=builder /app/converter .

EXPOSE 8080


CMD ["./converter"]