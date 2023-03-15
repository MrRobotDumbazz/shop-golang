FROM golang:1.19-alpine3.16 AS builder
LABEL stage=builder 
ENV GO111MODULE=on
WORKDIR /app 
COPY . .
RUN go install github.com/cosmtrek/air@latest
RUN apk add build-base && go build -o main ./cmd/main.go
# stage 2
FROM alpine:3.16 AS runner 
LABEL stage=runner 
LABEL maintainer="Made by Nurzhas && Mr.RobotDumbazz"
LABEL org.label-schema.description="Docker image for Shop"
WORKDIR /app
COPY --from=builder /app/main ./
COPY /internal /app/internal
COPY /static /app/static 
COPY /templates /app/templates
EXPOSE 8080 
RUN ls
CMD ["./main"]  