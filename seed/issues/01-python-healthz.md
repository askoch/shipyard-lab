# [python-api] Add a /healthz endpoint

Add `GET /healthz` to `python-api` that returns `200` and `{"status": "ok"}`.

Acceptance criteria:
- The route must not be shadowed by the `/{code}` catch-all route.
- Add a test in `python-api/tests/test_app.py`.
- Existing tests keep passing.
