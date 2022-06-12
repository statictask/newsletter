FROM golang:1.18-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /usr/src/app

RUN apk update \
  && apk add make git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ENV GOBUILDTAGS="migrate"

ADD . .
RUN make build


FROM alpine:3.16 AS final

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=builder /usr/src/app/build/newsletter /newsletter
COPY --from=builder /usr/src/app/build/migrate /migrate
COPY --from=builder /usr/src/app/db/migrations /db/migrations

EXPOSE 8080

ENTRYPOINT ["/newsletter"]

CMD ["server"]
