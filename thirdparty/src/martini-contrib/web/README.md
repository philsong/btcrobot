# web.Context
[hoisie][] [web.go][]'s Context for Martini.

[hoisie]:https://github.com/hoisie
[web.go]:https://github.com/hoisie/web

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/web)

## Description
`web.Context` provides a [web.go][] compitable layer for reusing the code written with
hoisie's `web.go` framework. Here compitable means we can use `web.Context` the same 
way as in hoisie's `web.go` but not the others.

## Usage

~~~ go
package main

import (
   "github.com/codegangsta/martini"
   "github.com/codegangsta/martini-contrib/web"
 )

func main() {
  m := martini.Classic()
  m.Use(web.ContextWithCookieSecret(""))

  m.Post("/hello", func(ctx *web.Context){
  	  ctx.WriteString("Hello World!")
  })

  m.Run()
}
~~~

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
* [Archs Sun](http://github.com/Archs)
* [hoisie][]
