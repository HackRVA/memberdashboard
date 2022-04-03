FROM golang:1.18 as backend-build

WORKDIR /membership
COPY . .

RUN go test -v ./...

RUN go mod vendor

RUN go build -o server -ldflags "-X memberserver/api.GitCommit=$(git rev-parse --short HEAD)" ./cmd/main

# create a file named Dockerfile
FROM node:14.17.6 as frontend-build

WORKDIR /app

COPY internal/web/package.json /app

RUN npm i

COPY ./web /app
# compile and bundle typescript
RUN npm run build

# copy from build environments
FROM node:14.17.6

WORKDIR /app

COPY --from=frontend-build /app/dist ./ui/dist/
COPY --from=backend-build /membership/server .
COPY --from=backend-build /membership/templates ./templates

ENTRYPOINT [ "./server" ]
