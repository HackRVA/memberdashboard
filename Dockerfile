FROM golang:1.15 as backend-build

WORKDIR /membership
COPY . .

RUN go mod vendor
RUN go build -o server

# create a file named Dockerfile
FROM node:latest as frontend-build

WORKDIR /app

COPY ui/package.json /app
RUN npm install


COPY ./ui /app
# compile and bundle typescript
RUN npm run rollup

# copy from build environments
FROM node:latest

WORKDIR /app

COPY --from=frontend-build /app/dist ./ui/dist/
COPY --from=backend-build /membership/server .


ENTRYPOINT [ "./server" ]
