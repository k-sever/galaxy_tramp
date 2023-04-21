FROM golang:1.20 AS build-env
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod\
    go mod download

FROM build-env AS build
RUN  --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
     go build -o /go/bin/game .

FROM build-env AS test
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    go test -v ./...

FROM gcr.io/distroless/static-debian11 AS game
WORKDIR /app
COPY --from=build /go/bin/game /game
ENTRYPOINT ["/game"]