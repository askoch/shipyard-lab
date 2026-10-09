# [go-api] Reject URLs longer than 2048 characters

`POST /links` stores URLs of any length. Reject URLs over 2048 characters with `400` and a JSON error body.

Acceptance criteria:
- A 2048 character URL is accepted; a 2049 character URL is rejected.
- Add tests for both boundaries.
