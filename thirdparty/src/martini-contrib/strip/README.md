# strip

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/strip)

## Description
packcage `strip` modifies the URL before the requests go into the other 
handlers.

Currently the main function in package strip is `strip.Prefix` which provides
the save functionality as `http.StripPrefix` and can be used in martini instance 
and request context level.

With `strip.Prefix` martini instances can be installed upon each other, and so
does some other web framework like [web.go][].

[web.go]:https://github.com/hoisie/web

## Usage

~~~ go
package main

import (
	"github.com/codegangsta/martini-contrib/strip"
	"github.com/codegangsta/martini"
)

func main() {
	m := martini.Classic()

	m2 := martini.Classic()
	m2.Get("/", func() string {
		return "Hello World from 2nd martini"
	})

	m2.Get("/foo", func() string {
		return "Hello foo"
	})

	m.Get("/", func() string {
		return "Hello World from 1st martini"
	})
	m.Get("/2ndMartini/.*", strip.Prefix("/2ndMartini"), m2.ServeHTTP)

	m.Run()
}
~~~

But the example above can only translate the same HTTP method from `m.Get`
to `m2.Get`, in order to transfer all kinds request such as `Post`,`Delete`,
etc to `m2`, martini has to provide a method `Any` to match any HTTP method
to a certain URL pattern.

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
* [Archs Sun](http://github.com/Archs)
* [hoisie](http://github.com/hoisie)
