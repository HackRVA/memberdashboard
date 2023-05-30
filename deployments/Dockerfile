FROM node:14.17.6 as frontend-build

WORKDIR /app

COPY ./web/package.json /app

RUN npm i

COPY ./web /app
# compile and bundle typescript
RUN npm run build

FROM golang:1.18 as backend-build

WORKDIR /membership
COPY . .
# copy ui files so that we can embed frontend into the binary
COPY --from=frontend-build /app/dist /membership/pkg/membermgr/ui/web

RUN go test -v ./...

RUN go mod vendor

RUN go build -o server -ldflags "-X memberserver/pkg/membermgr/controllers.GitCommit=$(git rev-parse --short HEAD)" ./cmd/main

# copy from build environments
FROM golang:1.18

WORKDIR /app

COPY --from=backend-build /membership/server .
COPY --from=backend-build /membership/pkg/membermgr/templates ./templates

ENTRYPOINT [ "./server" ]
