#!/bin/sh
# WynnMounts → WynnBreeder one-click migration script
# Run as root on the Alpine server.
set -e

OLD_NAME="wynnmounts"
NEW_NAME="wynnbreeder"
OLD_DIR="/opt/wynnmounts"
NEW_DIR="/opt/wynnbreeder"
REPO_DIR="$(cd "$(dirname "$0")" && pwd)"

# ── helpers ────────────────────────────────────────────────────────────────
info()  { printf '\033[1;34m==> %s\033[0m\n' "$*"; }
ok()    { printf '\033[1;32m    ✓ %s\033[0m\n' "$*"; }
warn()  { printf '\033[1;33m    ! %s\033[0m\n' "$*"; }
die()   { printf '\033[1;31mERROR: %s\033[0m\n' "$*" >&2; exit 1; }

# ── detect init system ─────────────────────────────────────────────────────
if command -v rc-service >/dev/null 2>&1; then
    INIT="openrc"
elif command -v systemctl >/dev/null 2>&1; then
    INIT="systemd"
else
    INIT="none"
fi

svc_stop()   {
    case "$INIT" in
        openrc)  rc-service "$1" stop  2>/dev/null || true ;;
        systemd) systemctl stop "$1"   2>/dev/null || true ;;
    esac
}
svc_disable() {
    case "$INIT" in
        openrc)  rc-update del "$1" default 2>/dev/null || true ;;
        systemd) systemctl disable "$1"     2>/dev/null || true ;;
    esac
}
svc_enable_start() {
    case "$INIT" in
        openrc)
            rc-update add "$1" default
            rc-service "$1" start
            ;;
        systemd)
            systemctl daemon-reload
            systemctl enable --now "$1"
            ;;
        none) warn "No init system detected — start the binary manually." ;;
    esac
}

# ── 1. stop & disable old service ─────────────────────────────────────────
info "Stopping old service ($OLD_NAME)..."
svc_stop    "$OLD_NAME"
svc_disable "$OLD_NAME"
ok "Service stopped"

# ── 2. create new directory ────────────────────────────────────────────────
info "Setting up $NEW_DIR..."
mkdir -p "$NEW_DIR/data"

# ── 3. migrate database ────────────────────────────────────────────────────
info "Migrating database..."

# Try common old DB locations
OLD_DB=""
for candidate in \
    "$OLD_DIR/data/${OLD_NAME}.db" \
    "$OLD_DIR/${OLD_NAME}.db" \
    "/data/${OLD_NAME}.db" \
    "./${OLD_NAME}.db"
do
    if [ -f "$candidate" ]; then
        OLD_DB="$candidate"
        break
    fi
done

NEW_DB="$NEW_DIR/data/${NEW_NAME}.db"

if [ -n "$OLD_DB" ]; then
    if [ -f "$NEW_DB" ]; then
        warn "New DB already exists at $NEW_DB — skipping copy (old: $OLD_DB)"
    else
        cp "$OLD_DB" "$NEW_DB"
        ok "Database copied: $OLD_DB → $NEW_DB"
    fi
else
    warn "Old database not found — starting fresh (checked $OLD_DIR/data/, $OLD_DIR/, /data/, ./)"
fi

# ── 4. build new binary ────────────────────────────────────────────────────
info "Building $NEW_NAME binary..."
if ! command -v go >/dev/null 2>&1; then
    die "Go not found. Install Go first: apk add go"
fi

cd "$REPO_DIR"
git pull origin main
go build -o "$NEW_DIR/$NEW_NAME" .
ok "Binary built: $NEW_DIR/$NEW_NAME"

# ── 5. install service ─────────────────────────────────────────────────────
info "Installing service ($INIT)..."

case "$INIT" in
    openrc)
        cat > "/etc/init.d/$NEW_NAME" <<EOF
#!/sbin/openrc-run
name="$NEW_NAME"
description="WynnBreeder"
command="$NEW_DIR/$NEW_NAME"
command_args="serve"
command_background=true
pidfile="/run/${NEW_NAME}.pid"
directory="$NEW_DIR"
environment="WYNNBREEDER_DB=$NEW_DB"

depend() {
    need net
}
EOF
        chmod +x "/etc/init.d/$NEW_NAME"
        ok "OpenRC service installed at /etc/init.d/$NEW_NAME"
        ;;
    systemd)
        cat > "/etc/systemd/system/${NEW_NAME}.service" <<EOF
[Unit]
Description=WynnBreeder
After=network.target

[Service]
ExecStart=$NEW_DIR/$NEW_NAME serve
WorkingDirectory=$NEW_DIR
Restart=always
Environment=WYNNBREEDER_DB=$NEW_DB

[Install]
WantedBy=multi-user.target
EOF
        ok "systemd unit installed at /etc/systemd/system/${NEW_NAME}.service"
        ;;
    none)
        warn "No init system — skipping service install"
        ;;
esac

# ── 6. enable & start ─────────────────────────────────────────────────────
info "Starting $NEW_NAME..."
svc_enable_start "$NEW_NAME"

# ── 7. summary ────────────────────────────────────────────────────────────
printf '\n\033[1;32mMigration complete!\033[0m\n'
printf '  Binary : %s\n' "$NEW_DIR/$NEW_NAME"
printf '  DB     : %s\n' "$NEW_DB"
printf '  Service: %s (%s)\n' "$NEW_NAME" "$INIT"
printf '\nOld directory (%s) left intact — remove manually once verified:\n' "$OLD_DIR"
printf '  rm -rf %s\n' "$OLD_DIR"
