FROM golang:1.21 as build
WORKDIR /var/task

RUN apt update \
 && apt install -y \
      imagemagick \
      libmagickwand-dev

COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN GO111MODULE=on GOOS=linux go build -ldflags="-s -w" -o build/api        pkg/handlers/api/lambda/main.go \
 && GO111MODULE=on GOOS=linux go build -ldflags="-s -w" -o build/deletelgtm pkg/handlers/deletelgtm/main.go

ENTRYPOINT []

# ---------------------------

FROM golang:1.21
WORKDIR /var/task

RUN apt update \
 && apt install -y \
      imagemagick \
      libmagickwand-dev
COPY --from=build /var/task/pkg/static/ /var/task/pkg/static/
COPY --from=build /var/task/build/api /var/task/build/api
COPY --from=build /var/task/build/deletelgtm /var/task/build/deletelgtm

ENTRYPOINT []
