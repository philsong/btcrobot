# secure
Martini middleware that helps enable some quick security wins.

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/secure)

## Usage

~~~ go
import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/secure"
)

func main() {
  m := martini.Classic()

  martini.Env = martini.Prod  // You have to set the environment to `production` for all of secure to work properly!

  m.Use(secure.Secure(secure.Options{
    AllowedHosts: []string{"example.com", "ssl.example.com"},
    SSLRedirect: true,
    SSLHost: "ssl.example.com",
    SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
    STSSeconds: 315360000,
    STSIncludeSubdomains: true,
    FrameDeny: true,
    ContentTypeNosniff: true,
    BrowserXssFilter: true,
    ContentSecurityPolicy: "default-src 'self'",
  }))
  m.Run()
}

~~~

Make sure to include the secure middleware as close to the top as possible. It's best to do the allowed hosts and SSL check before anything else.

The above example will only allow requests with a host name of 'example.com', or 'ssl.example.com'. Also if the request is not https, it will be redirected to https with the host name of 'ssl.example.com'.
After this it will add the following headers:
~~~
Strict-Transport-Security: 315360000; includeSubdomains
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
~~~

###Set the `MARTINI_ENV` environment variable to `production` when deploying!
If you don't, the SSLRedirect and STS Header will not work. This allows you to work in development/test mode and not have any annoying redirects to HTTPS (ie. development can happen on http). If this is not the behavior you're expecting, see the `DisableProdCheck` below in the options.

You can also disable the production check for testing like so:
~~~ go
//...
m.Use(secure.Secure(secure.Options{
    SSLRedirect: true,
    STSSeconds: 315360000,
    DisableProdCheck: martini.Env == martini.Test,
  }))
//...
~~~


### Options
`secure.Secure` comes with a variety of configuration options:

~~~ go
// ...
m.Use(secure.Secure(secure.Secure{
  AllowedHosts: []string{"ssl.example.com"}, // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
  SSLRedirect: true, // If SSLRedirect is set to true, then only allow https requests. Default is false.
  SSLHost: "ssl.example.com", // SSLHost is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host.
  SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid https request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
  STSSeconds: 315360000, // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
  STSIncludeSubdomains: true, // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
  FrameDeny: true, // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
  CustomFrameOptionsValue: "SAMEORIGIN", // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option.
  ContentTypeNosniff: true, // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
  BrowserXssFilter: true, // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
  ContentSecurityPolicy: "default-src 'self'", // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
  DisableProdCheck: true, // This will ignore our production check, and will follow the SSLRedirect and STSSeconds/STSIncludeSubdomains options... even in development! This would likely only be used to mimic a production environment on your local development machine.
}))
// ...
~~~

### Nginx
If you would like to add the above security rules directly to your nginx configuration, everything is below:
~~~
# Allowed Hosts:
if ($host !~* ^(example.com|ssl.example.com)$ ) {
    return 500;
}

# SSL Redirect:
server {
    listen      80;
    server_name example.com ssl.example.com;
    return 301 https://ssl.example.com$request_uri;
}

# Headers to be added:
add_header Strict-Transport-Security "max-age=315360000";
add_header X-Frame-Options "DENY";
add_header X-Content-Type-Options "nosniff";
add_header X-XSS-Protection "1; mode=block";
add_header Content-Security-Policy "default-src 'self'";
~~~

## Authors
* [Cory Jacobsen](http://github.com/cojac)
