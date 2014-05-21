# cors

Martini middleware/handler to enable CORS support.

## Usage

~~~ go
import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/cors"
)

func main() {
  m := martini.Classic()
  // CORS for https://foo.* origins, allowing:
  // - PUT and PATCH methods
  // - Origin header
  // - Credentials share
  m.Use(cors.Allow(&cors.Options{
    AllowOrigins:     []string{"https://foo\\.*"},
    AllowMethods:     []string{"PUT", "PATCH"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
  }))
  m.Run()
}
~~~

## Authors

* [Burcu Dogan](http://github.com/rakyll)
