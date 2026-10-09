# [python-api] Validate URLs on link creation

`POST /links` accepts any string as the URL. Only accept `http` and `https` URLs with a host.

Acceptance criteria:
- `javascript:alert(1)`, empty strings, and strings without a host return `422`.
- Valid `http` and `https` URLs still return `201`.
- Add tests for both the rejected and accepted cases.
