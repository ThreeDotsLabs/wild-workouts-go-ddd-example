#!/usr/bin/env bash
set -euo pipefail

# Fix Docker socket permissions by matching the container's docker group GID
# to the GID that owns the mounted socket.
if [ -S /var/run/docker.sock ]; then
    SOCK_GID=$(stat -c '%g' /var/run/docker.sock)
    if ! getent group "$SOCK_GID" > /dev/null 2>&1; then
        addgroup -g "$SOCK_GID" -S dockerhost
        addgroup claude dockerhost
    else
        GROUP_NAME=$(getent group "$SOCK_GID" | cut -d: -f1)
        addgroup claude "$GROUP_NAME" 2>/dev/null || true
    fi
fi

# Seed .claude/ from host on first run (volume empty).
# Once seeded, the persistent volume owns the data.
if [ -d /home/claude/.claude-host ] && [ ! -f /home/claude/.claude/settings.json ]; then
    cp -a /home/claude/.claude-host/. /home/claude/.claude/
fi
chown -R claude:claude /home/claude/.claude

# Seed .claude.json from host on every start.
# .claude.json is NOT in a named volume, so it always resets to the image layer on each run.
# Always overwrite with the host config so Claude Code starts with the correct settings.
if [ -f /home/claude/.claude.json.host ]; then
    cp /home/claude/.claude.json.host /home/claude/.claude.json
    chown claude:claude /home/claude/.claude.json
fi

# Fix npm cache volume permissions (created as root by Docker).
echo "Setting up npm cache permissions..."
chown -R claude:claude /home/claude/.npm 2>/dev/null || true

# Fix git worktree paths for Docker.
# Worktree .git files contain absolute host paths that don't exist in the container.
# When the main repo's .git is mounted at /repo-git, redirect git via env vars.
if [ -f /src/.git ] && [ -d /repo-git ]; then
    GITDIR_HOST=$(sed 's/^gitdir: //' /src/.git)
    WORKTREE_NAME=$(basename "$GITDIR_HOST")
    export GIT_DIR="/repo-git/worktrees/$WORKTREE_NAME"
    export GIT_WORK_TREE="/src"
fi

# Suppress the "Run /terminal-setup" startup tip — it doesn't apply inside Docker.
# Setting a far-future "last shown" session count keeps the tip permanently in cooldown
# without falsely claiming Shift+Enter is installed (which would change the toolbar hint).
if [ -f /home/claude/.claude.json ]; then
    jq '.tipsHistory["terminal-setup"] = 99999999' /home/claude/.claude.json \
        > /tmp/.claude.json.tmp \
        && mv /tmp/.claude.json.tmp /home/claude/.claude.json \
        && chown claude:claude /home/claude/.claude.json
else
    printf '{"tipsHistory":{"terminal-setup":99999999}}\n' > /home/claude/.claude.json
    chown claude:claude /home/claude/.claude.json
fi

# Unset multiplexer env vars — these are inherited from the host shell but don't
# apply inside the container. TERM_PROGRAM is kept so Claude Code can identify the
# terminal emulator and enable Shift+Enter support (kitty/iTerm2 keyboard protocol).
unset TMUX STY

echo "Starting Claude Code..."
exec su-exec claude claude "$@"
