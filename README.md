UPDATE:

The answer is here: https://stackoverflow.com/questions/43166750/goroutines-with-sync-waitgroup-end-before-last-wg-done

What you read below is now historical context for the original problem I originally faced (commit df727edab93fff2ff07f0176abeecb50925a1ad6).  It's now fixed (commit a68447af6217d03539d4db2a4ca963b0147807d3) The world makes sense again.

```
go mod tidy
```
Then
```
go test -v -race -count=1 ./implementations/...
```

Rerun a few times.  Sometimes the output looks good, where `after wg.Wait()` shows a fully accumulated slice:
```
=== RUN   TestIt
MIKE: got: thing1
MIKE: size: 1
MIKE: got: thing2
MIKE: size: 2
MIKE: got: thing3
MIKE: size: 3
MIKE: after wg.Wait(): 3
--- PASS: TestIt (0.00s)
PASS
```

Sometimes is looks bad, where `after wg.Wait()` (apparently?) shows wg.Wait() released too soon.  I don't know...
```
=== RUN   TestIt
MIKE: got: thing1
MIKE: size: 1
MIKE: got: thing2
MIKE: size: 2
MIKE: got: thing3
MIKE: after wg.Wait(): 2
MIKE: size: 3
--- PASS: TestIt (0.00s)
PASS
```
