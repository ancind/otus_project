FROM golang:1.18-alpine AS build

WORKDIR $GOPATH/src/github.com/ancind/otus-project

COPY . ./
RUN go mod download
RUN go build -o /go/bin/image-previewer ./cmd/server/main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/image-previewer /bin/

ENTRYPOINT ["image-previewer"]