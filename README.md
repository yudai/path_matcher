# Path Matcher

## Usage

```sh
go get github.com:yudai/pmatcher
```

```go
import "github.com:yudai/pmatcher/"

matcher := pmatcher.New()
matcher.Add("/foo/bar")
matcher.Add("/foo/:p1/:p2")
matcher.Add("/foo/bar/baz")
matcher.Add("/baz/:p1")
matcher.Add("/baz/one")
matcher.Add("/baz/:p1/two")
matcher.Add("/baz/:p2/three")

matched, pattern, params := matcher.Match("/foo/some/thing")
// matched => true
// pattern => "/foo/:p1/:p2"
// params  => map[string]string{"p1": "some", "p2": "thing"}
```
