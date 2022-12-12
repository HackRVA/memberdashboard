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

RUN go test -v ./...

RUN go mod vendor

# copy ui files so that we can embed frontend into the binary
COPY --from=frontend-build /app/dist ./internal/ui/web/

RUN go build -o server -ldflags "-X memberserver/internal/controllers.GitCommit=$(git rev-parse --short HEAD)" ./cmd/main

# copy from build environments
FROM golang:1.18

WORKDIR /app

COPY --from=backend-build /membership/server .
COPY --from=backend-build /membership/internal/templates ./templates

ENTRYPOINT [ "./server" ]
