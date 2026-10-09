# [node-api] Click counts in /links/:code/stats are off by one

After following a short link once, `GET /links/:code/stats` reports `0` clicks. A link that has never been followed reports `-1`. The Python and Go services report the correct number.

Acceptance criteria:
- Stats return the true click count, never negative.
- Add a regression test that follows a link twice and asserts `clicks === 2`.
