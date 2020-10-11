FROM golang:1.11.0
ADD AnimeServer AnimeServer
EXPOSE 8060
LABEL name=AnimeServer
ENTRYPOINT ["./AnimeServer"]