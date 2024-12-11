FROM golang:1.22-alpine as build

WORKDIR /app

# COPY /root/* /root/

ENV GOPRIVATE=github.com/be-growth

RUN apk update && apk upgrade && apk add git

RUN git config --global url."https://"$( cat /root/token )":x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build main.go

FROM alpine:3.14
WORKDIR /app

COPY --from=build /app/main /app/main

EXPOSE ${PORT}

CMD [ "./main" ]