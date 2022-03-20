FROM golang:1.15-alpine AS build
WORKDIR /home/
#RUN mkdir gannant-code-test
#ADD gannant-code-test gannant-code-test
COPY . ./gannant-code-test/server.exe
WORKDIR gannant-code-test
RUN apk add make
RUN apk add git
RUN apk add gcc
RUN apk add build-base
ENTRYPOINT ["make", "all"]