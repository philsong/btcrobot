# acceptlang
Using the `acceptlang` handler you can automatically parse the `Accept-Language` HTTP header and expose it as an `AcceptLanguages` slice in your handler functions. The `AcceptLanguages` slice contains `AcceptLanguage` values, each of which represent a qualified (or unqualified) language. The values in the slice are sorted descending by qualification (the most qualified languages will have the lowest indexes in the slice).

Unqualified languages are interpreted as having the maximum qualification of `1`, as defined in the HTTP/1.1 specification.

For more information:
* [HTTP/1.1 Accept-Language specification](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) 
* [API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/acceptlang)

## Usage
Simply add a new handler function instance to your handler chain using the `acceptlang.Languages()` function as well as an `AcceptLanguages` dependency in your handler function. The `AcceptLanguages` dependency will be satisified by the handler.

For example:

```go
package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/acceptlang"
	"net/http"
)

func main() {
	m := martini.Classic()

	m.Get("/", acceptlang.Languages(), func(languages acceptlang.AcceptLanguages) string {
		return fmt.Sprintf("Languages: %s", languages)
	})

	http.ListenAndServe(":8090", m)
}
```

## Authors
* [Tom Bruggeman](http://github.com/tmbrggmn)
