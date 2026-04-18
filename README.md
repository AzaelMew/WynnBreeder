# WynnBreeder

A self-hosted web app for tracking Wynncraft mount breeding data. Submit Parent A + Parent B + Offspring stats, view stat inheritance patterns, and analyze breeding outcomes over time.

## Features

- Import mount stats directly from clipboard (copy JSON from the game)
- Save a pending breeding with just the parents — come back later to add the offspring
- Per-user ownership: only you can complete your own pending submissions
- Admin panel for user and submission management
- Analytics: stat inheritance, potential trends, color frequency tables
- Single binary, no external dependencies, SQLite database

## Requirements

- Go 1.21+ (for self-build) **or** Docker

## Setup

### Option A — Go binary

**1. Build**
```bash
go build -o wynnbreeder .
```

**2. Create the admin account** (first time only)
```bash
./wynnbreeder seed-admin --username admin --password yourpassword
```

**3. Start the server**
```bash
./wynnbreeder serve --port 8080
```

Open `http://localhost:8080` in your browser and log in.

---

### Option B — Docker Compose

```bash
docker compose up -d
```

The database is stored in `./data/wynnbreeder.db` on your host.

Create the admin account inside the container:
```bash
docker compose exec wynnbreeder ./wynnbreeder seed-admin --username admin --password yourpassword
```

---

## Configuration

All settings can be set via environment variables or CLI flags.

| Environment variable | CLI flag | Default | Description |
|---|---|---|---|
| `WYNNBREEDER_PORT` | `--port` | `8080` | Port to listen on |
| `WYNNBREEDER_DB` | `--db` | `./wynnbreeder.db` | Path to SQLite database file |
| `WYNNBREEDER_SESSION_DAYS` | — | `30` | Session cookie lifetime in days |

---

## User management

Only the admin can create accounts — there is no self-registration.

1. Log in as admin and go to `/admin`
2. Click **New User**, fill in a username and password
3. Share the credentials with the person you want to give access to

Users can submit breeding data. Admins can delete any submission and manage all users.

---

## Submitting breeding data

1. Go to `/submit`
2. Copy the mount JSON from Wynncraft and paste it into the **Parent A** slot
3. Do the same for **Parent B**
4. If breeding is done, paste the **Offspring** too and click **Submit Complete Breeding**
5. If the offspring isn't born yet, click **Save Parents (breed in progress)** — you'll see it on the dashboard and can come back later to add the offspring

---

## Database size

Each breeding record (3 mounts) uses approximately **1 KB** of disk space.

| Submissions | Approximate size |
|---|---|
| 1,000 | ~1 MB |
| 10,000 | ~10 MB |
| 100,000 | ~100 MB |

**Recommended disk allocation: 5–10 GB**, which is far more than will ever be needed for the database. The Go binary itself is ~15 MB.

---

## Running on a server

The app is very lightweight. Minimum recommended specs:

- **CPU:** 1 vCPU (shared is fine)
- **RAM:** 512 MB (2 GB is comfortable headroom)
- **Disk:** 5–10 GB

A basic VPS or a small cloud instance is more than sufficient.

### Keeping it running with systemd

Create `/etc/systemd/system/wynnbreeder.service`:

```ini
[Unit]
Description=WynnBreeder
After=network.target

[Service]
ExecStart=/opt/wynnbreeder/wynnbreeder serve
WorkingDirectory=/opt/wynnbreeder
Restart=always
Environment=WYNNBREEDER_DB=/opt/wynnbreeder/data/wynnbreeder.db

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable --now wynnbreeder
```

---

## Backup

The entire database is a single file. Back it up by copying it:

```bash
cp wynnbreeder.db wynnbreeder.db.bak
```

For automated backups, SQLite's online backup is safe to run while the server is live:

```bash
sqlite3 wynnbreeder.db ".backup wynnbreeder-backup.db"
```
