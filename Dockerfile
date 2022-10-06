FROM golang:alpine as build
WORKDIR /src
COPY /src .
RUN go version
RUN go build -o server main.go

FROM alpine
COPY --from=build /src/server server
COPY --from=build /src/index.html index.html
ENTRYPOINT ["./server"]