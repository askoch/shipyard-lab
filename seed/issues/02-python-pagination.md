# [python-api] Paginate GET /links

`GET /links` currently returns every link. Add `limit` (default 20, max 100) and `offset` (default 0) query parameters.

Acceptance criteria:
- Invalid values (negative offset, limit over 100) return `422`.
- Response shape stays a JSON array.
- Add tests covering the default, a custom limit, and an out-of-range limit.
