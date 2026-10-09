# [go-api] Add a /healthz endpoint

Add `GET /healthz` to `go-api` returning `200` and `{"status":"ok"}`.

Acceptance criteria:
- The route is registered so it is not caught by `GET /{code}`.
- Add a test in `go-api/main_test.go`.
