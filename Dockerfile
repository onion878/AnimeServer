FROM golang:1.11.0
ADD main main
EXPOSE 8000
LABEL name=AnimeServer
ENTRYPOINT ["./main"]