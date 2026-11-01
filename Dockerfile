# Build FrontEnd
FROM node:26-alpine3.24 AS frontend
RUN apk upgrade --no-cache
WORKDIR /frontend/
COPY ./frontend/src ./src
COPY ./frontend/public ./public
COPY ./frontend/favicon.svg ./favicon.svg
COPY ./frontend/README.md ./README.md
COPY ./frontend/index.html ./index.html
COPY ./frontend/package-lock.json ./package-lock.json
COPY ./frontend/.eslintrc ./.eslintrc
COPY ./frontend/jsconfig.json ./jsconfig.json
COPY ./frontend/package.json ./package.json
COPY ./frontend/vite.config.mjs ./vite.config.mjs
RUN npm ci
RUN npm run build

# Build
FROM golang:1.27.2-alpine3.24 AS build
RUN apk upgrade --no-cache
WORKDIR /build/
COPY go.mod .
RUN go mod tidy
COPY . .
COPY --from=frontend /frontend/dist/assets ./assets
COPY --from=frontend ./frontend/dist/index.html ./
COPY --from=frontend /frontend/dist/*.webmanifest /frontend/dist/*.js /frontend/dist/*.png ./web/
RUN go version
RUN go build

# Production
FROM alpine:3.24 AS production
RUN apk upgrade --no-cache
WORKDIR /monitoruimanager/

RUN apk upgrade --no-cache && apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=build /build/monitoruimanager ./
COPY ./*.crt ./
COPY ./*.key ./
COPY ./*.pem ./
COPY ./monitoruimanager.json ./

EXPOSE 7007

CMD ["./monitoruimanager", "-log=stdout"]
