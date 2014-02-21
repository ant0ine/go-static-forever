
Go-Static-Forever
=================

*Serve files that never change*

[![Build Status](https://travis-ci.org/ant0ine/go-static-forever.png?branch=master)](https://travis-ci.org/ant0ine/go-static-forever)

[![Total views](https://sourcegraph.com/api/repos/github.com/ant0ine/go-static-forever/counters/views.png)](https://sourcegraph.com/github.com/ant0ine/go-static-forever)

If you can put a version string (commitish for instance) in the path of your
static files, then the content served by the corresponding URLs is guaranteed
to never change. A whole set of optimisations become possible.

* If the request contains `If-Modified-Since`, return `304` without checking anything

* Set the `Expires` to `<forever>` (`<forever>` defaulting to one year)

* Set the `Cache-Control` header to `public; max-age=<forever>; s-maxage=<forever>`

* Set the `Last-Modified` headers to `<origin>` (`<origin>` being 1970)

This handler is implemented as a wrapper around http.FileServer, and when the
isDevelopment flag is set, http.FileServer is used directly.

Install
-------

This package is "go-gettable", just do:

    go get github.com/ant0ine/go-static-forever

Example
-------

~~~ go
package main

import(
        "github.com/ant0ine/go-static-forever"
        "net/http"
)

func main() {
        handler := forever.NewStaticHandler(
                http.Dir("/static/"),   // FileSytem to serve
                "1234567",              // version string, like a commitish for instance
                nil,                    // "forever duration", defaults to one year
                false,                  // isDevelopement
        )

        http.ListenAndServe(":8080", handler)
}
~~~

Documentation
-------------

- [Online Documentation (godoc.org)](http://godoc.org/github.com/ant0ine/go-static-forever)

Copyright (c) 2013 Antoine Imbert

[MIT License](https://github.com/ant0ine/go-static-forever/blob/master/LICENSE)


[![Analytics](https://ga-beacon.appspot.com/UA-309210-4/go-static-forever/readme)](https://github.com/igrigorik/ga-beacon)
