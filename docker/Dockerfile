FROM golang:alpine AS build-env

RUN apk add --no-cach git
RUN adduser -D -u 10000 gosrvusr
RUN mkdir /bank-go/ && chown gosrvusr /bank-go/

USER gosrvusr
WORKDIR /bank-go/
ADD . /bank-go/

RUN CGO_ENABLED=0 go build -o /bank-go/app .

FROM alpine:latest
RUN adduser -D -u 10000 gosrvusr
USER gosrvusr

WORKDIR /

COPY --from=build-env /bank-go/app /
EXPOSE 8080

CMD ["/app"]