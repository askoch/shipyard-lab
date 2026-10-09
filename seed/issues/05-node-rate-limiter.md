# [node-api] Add a simple per-IP rate limiter

Limit each client IP to 60 requests per minute across all routes. Over the limit, return `429` with a `Retry-After` header.

Acceptance criteria:
- No new runtime dependencies; an in-memory fixed-window counter is fine.
- The limit is configurable through a `RATE_LIMIT_PER_MIN` environment variable.
- Add a test that exceeds a low limit and sees a `429`.
