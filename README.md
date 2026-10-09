# shipyard-lab

A small sandbox for trying Upsun Dispatch. It contains one tiny service, a link shortener with click counts, implemented three times so you can compare how each runtime builds, tests, and deploys.

| Service | Runtime | Notes |
| --- | --- | --- |
| `go-api` | Go 1.22+ | `net/http` only, no dependencies |
| `node-api` | Node 20+ | Built-in `http` module, no dependencies |
| `python-api` | Python 3.12+ | FastAPI and uvicorn |

All three expose the same API with in-memory storage:

- `POST /links` with `{"url": "..."}` returns `201` and the new link
- `GET /links` lists links
- `GET /{code}` redirects (`307`) and counts a click
- `GET /links/{code}/stats` returns the link and its click count

## Run and test locally

```
cd go-api && go test ./... && go run .
cd node-api && npm test && npm start
cd python-api && pip install -r requirements.txt && python -m pytest -q && uvicorn app.main:app
```

CI (`.github/workflows/ci.yml`) runs all three test suites on every pull request.

## Deploying on Upsun

`.upsun/config.yaml` defines all three apps in one project, each with its own route. Treat it as a starting point and validate it with the Upsun CLI before the first deploy, since runtime version strings change over time.

## Using it with Dispatch

Connect the repo to Dispatch, then run its two launch workflows.

### Issue to Code

`seed/issues/` holds eight small, well-scoped issues with acceptance criteria, spread across the three services. `scripts/seed.sh` files them for you. Good first picks:

1. `[python-api] Add a /healthz endpoint` (easy win, shows the loop end to end)
2. `[node-api] Click counts in stats are off by one` (bug fix with a regression test)
3. `[python-api] Paginate GET /links` (touches validation and tests)

### Automated Code Review

`seed/patches/` holds three pull requests with deliberate flaws. `scripts/seed.sh` opens them. Do not fix them before the review runs. What to look for:

| Patch | Planted flaws |
| --- | --- |
| `01-python-click-log` | SQL injection through string-built INSERT using the path and the User-Agent header; no test; a new SQLite file written to disk |
| `02-node-admin-delete` | Hardcoded admin key in source; loose equality on the key check; no test; deletes without checking the link exists |
| `03-go-top-links` | Data race (map read without the mutex); ignored `Atoi` error; unchecked slice bound that panics when `n` exceeds the number of links; no test |

Also present in the baseline: the Node stats endpoint reports `clicks - 1`, and the existing Node test does not assert the count. That is the bug behind issue 4.

### Suggested demo script (3 minutes)

1. Show the open issue and its acceptance criteria.
2. Run the Issue-to-Code workflow and point out the human gate before merge.
3. Open the resulting pull request and show the audit trail (prompts, diffs, approvals) and the cost of the run.
4. Run Automated Code Review on a seeded pull request and compare its findings to the table above.

## Seeding

```
gh auth login
gh repo create shipyard-lab --public --source . --push   # first time only
./scripts/seed.sh
```

Run `seed.sh` once. Running it again creates duplicate issues and pull requests.
