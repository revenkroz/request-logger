FROM node:20-alpine AS frontend

WORKDIR /app/frontend

# cache dependencies
RUN --mount=type=bind,source=frontend/package.json,target=package.json \
    --mount=type=bind,source=frontend/yarn.lock,target=yarn.lock \
    --mount=type=cache,target=/root/.yarn \
    yarn install --frozen-lockfile

# build
COPY frontend .

RUN yarn build


FROM golang:1.22-alpine AS backend

RUN apk --no-cache add bash git openssh

WORKDIR /app

# cache dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

# build
COPY . .
COPY --from=frontend /app/frontend/dist/ ./frontend/dist/

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o requestlogger -ldflags "-w -s" .


FROM alpine:3.16 AS prod

WORKDIR /app

COPY --from=backend /app/requestlogger /app/

CMD [ "/app/requestlogger" ]
