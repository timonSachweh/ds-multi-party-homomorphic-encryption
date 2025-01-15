FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ARG TARGETOS TARGETARCH

COPY . .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/application ./cmd/main.go

FROM alpine

RUN apk --no-cache add curl

COPY --from=build /out/application /app/application

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT [ "/app/application" ]