FROM node:18.20.4 as frontend-build

WORKDIR /app

COPY ./web/package.json /app

RUN npm i

COPY ./web /app
# compile and bundle typescript
RUN npm run build

FROM golang:1.22 as backend-build

WORKDIR /membership
COPY . .
# copy ui files so that we can embed frontend into the binary
COPY --from=frontend-build /app/dist/ /membership/web/dist

RUN go test -v ./...

RUN go mod vendor

RUN go build -o server -ldflags "-X github.com/HackRVA/memberserver/controllers.GitCommit=$(git rev-parse --short HEAD)" ./cmd/main

# copy from build environments
FROM golang:1.22

WORKDIR /app

COPY --from=backend-build /membership/server .
COPY --from=backend-build /membership/templates ./templates

ENTRYPOINT [ "./server" ]
