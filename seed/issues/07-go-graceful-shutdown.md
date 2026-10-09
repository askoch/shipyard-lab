# [go-api] Shut down gracefully on SIGTERM

`main()` uses `http.ListenAndServe` directly, so in-flight requests are cut off on deploy. Switch to an `http.Server` and shut it down gracefully on `SIGINT` and `SIGTERM` with a 10 second timeout.

Acceptance criteria:
- Behavior of all existing routes is unchanged.
- The shutdown logic is small and readable; extract it into a helper if it helps testing.
