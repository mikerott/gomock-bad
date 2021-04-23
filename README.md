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
