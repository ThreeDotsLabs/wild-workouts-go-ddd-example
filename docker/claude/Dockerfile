FROM node:22-alpine

RUN apk add --no-cache \
    git bash docker-cli curl jq github-cli zsh su-exec \
    libgcc libstdc++ ripgrep poppler-utils \
    chromium \
    go go-task

RUN addgroup -S docker 2>/dev/null || true \
    && adduser -D -s /bin/zsh claude \
    && addgroup claude docker

ENV GOPATH=/home/claude/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p /home/claude/go && chown claude:claude /home/claude/go

RUN npm install -g @anthropic-ai/claude-code@latest

ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PLAYWRIGHT_MCP_EXECUTABLE_PATH=/usr/bin/chromium-browser

RUN ln -s /usr/bin/sed /usr/local/bin/gsed \
    && ln -s /usr/bin/sort /usr/local/bin/gsort \
    && ln -s /usr/bin/uniq /usr/local/bin/guniq

ENV USE_BUILTIN_RIPGREP=0

WORKDIR /src

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
