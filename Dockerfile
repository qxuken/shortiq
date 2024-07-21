FROM node:lts-alpine as frontend-builder

WORKDIR /app

COPY . .

RUN npm ci
RUN npm run build

# ---

FROM golang:1.22-alpine as backend-builder

WORKDIR /app

RUN apk add build-base

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/short_web

# ---

FROM alpine:latest

RUN mkdir /data
VOLUME /data

WORKDIR /app

COPY --from=backend-builder /app/short_web /app/short_web
COPY --from=frontend-builder /app/assets/ /app/assets/

ENV SHORTIQ_DEBUG="false"
ENV SHORTIQ_DATA_PATH="/data"
ENV SHORTIQ_BIND="0.0.0.0"
ENV SHORTIQ_PORT="8080"
ENV SHORTIQ_PUBLIC_URL="http://localhost:8080"

EXPOSE 8080
CMD ["/app/short_web"]
