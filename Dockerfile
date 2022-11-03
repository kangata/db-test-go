FROM golang:1.19-alpine as builder

WORKDIR /build

COPY . .

RUN go build -o dbtest

FROM alpine

RUN apk add --no-cache tzdata

COPY --from=builder /build/dbtest /bin/dbtest

CMD ["/bin/dbtest"]
