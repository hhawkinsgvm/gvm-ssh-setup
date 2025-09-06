# --- build stage ---
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/gvm-ssh ./main.go

# --- runtime ---
FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    openssh-client git ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# Install glab CLI
RUN curl -fsSL https://gitlab.com/gitlab-org/cli/-/releases/v1.36.0/downloads/glab_1.36.0_Linux_x86_64.tar.gz | \
    tar -xzC /usr/local/bin --strip-components=1 glab_1.36.0_Linux_x86_64/bin/glab && \
    chmod +x /usr/local/bin/glab

COPY --from=build /out/gvm-ssh /usr/local/bin/gvm-ssh
COPY entrypoint.sh /usr/local/bin/entrypoint
RUN chmod +x /usr/local/bin/entrypoint
ENTRYPOINT ["entrypoint"]
CMD ["wizard"]