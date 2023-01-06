# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS builder
WORKDIR /forum

COPY . .

RUN go mod tidy
RUN apk add build-base && go build  cmd/main.go


FROM alpine
WORKDIR /forum
COPY --from=builder /forum .
LABEL maintainers = "HgCl2 && aleebeg && Alfarabi09"
LABEL version = "1.0.0"
EXPOSE 4888
CMD ["/forum/main"]
