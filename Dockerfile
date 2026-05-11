FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

ARG GITLAB_ACCESS_TOKEN
ARG GITLAB_DOMAIN=gitlab.com

# .netrc
RUN echo "machine ${GITLAB_DOMAIN}" > ~/.netrc && \
    echo "login oauth2" >> ~/.netrc && \
    echo "password ${GITLAB_ACCESS_TOKEN}" >> ~/.netrc && \
    chmod 600 ~/.netrc

RUN git config --global url."https://${GITLAB_DOMAIN}/".insteadOf "https://${GITLAB_DOMAIN}/"

WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GOPRIVATE=${GITLAB_DOMAIN}/_spacemc_/* && \
    go mod download -x

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g internal/infra/http/http.go --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o users-service ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

COPY --from=builder /app/users-service /app/users-service

COPY --from=builder /app/docs /app/docs

ENV SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080 \
    ENV=production

EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/users-service"]