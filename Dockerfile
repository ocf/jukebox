# build stage
FROM golang:1.18 AS build-env
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build .

# final stage
FROM python:3-slim-bullseye
RUN ["pip3", "install", "yt-dlp"]
WORKDIR /app
COPY dist ./dist
COPY --from=build-env /src/Jukebox /app/
ENTRYPOINT [ "./Jukebox" ]
