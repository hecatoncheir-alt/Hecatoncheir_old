FROM golang
MAINTAINER Vostrikov Vitaliy
WORKDIR src/github.com/hecatoncheir/Hecatoncheir

COPY . .
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 8080
CMD ["go-wrapper", "run", "main.go"]