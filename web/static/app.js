// WynnBreeder — global JS utilities
// Alpine.js components are inline in templates

// CSRF: not needed for cookie-based auth with SameSite=Lax
// Fetch wrapper for consistent error handling
async function apiFetch(url, opts = {}) {
    const res = await fetch(url, {
        headers: { 'Content-Type': 'application/json', ...opts.headers },
        ...opts,
    });
    const data = await res.json().catch(() => ({}));
    return { ok: res.ok, status: res.status, data };
}
