FROM golang:1.18-alpine as builder

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

RUN apk --no-cache add ca-certificates
COPY . /src
WORKDIR /src

RUN go build -o /sc-bot .
RUN mkdir -p "/dumps"

FROM scratch

COPY --from=builder /sc-bot /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dumps /

CMD ["/sc-bot"]
