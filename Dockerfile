FROM golang:1.25-alpine AS build

RUN apk add --no-cache git

ARG GITLAB_TOKEN
RUN echo "Token length: $(echo ${#GITLAB_TOKEN})" && \
    if [ -z "$GITLAB_TOKEN" ]; then \
        echo "ERROR: GITLAB_TOKEN is empty!"; \
        exit 1; \
    fi

RUN git config --global url."https://gitlab-ci-token:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/"

RUN git config --global --list

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/main ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=build /app/main .
EXPOSE 8081
CMD ["./main"]