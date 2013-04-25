
Go-Static-Forever
=================

*Serve files that never change*

If you can put a version string (commitish for instance) in the path of your
static files, then the content served by the corresponding URLs is guaranteed
to never change. A whole set of optimisations become possible.

* Set the Expires and Last-Modified headers to <forever>
* Set the Cache-Control header to "public; max-age=<forever>; s-maxage=<forever>"
* Set the Etag header to the full file path ? TODO
* If the request contains If-Modified-Since, return 304 without checking anything

This handler is implemented as a wrapper around http.FileServer, and when the
isDevelopment flag is set, http.FileServer is used directly.

Install
-------

This package is "go-gettable", just do:

    go get github.com/ant0ine/go-static-forever

Example
-------

    package main

    import(
            "github.com/ant0ine/go-static-forever"
            "net/http"
    )

    handler := forever.NewStaticHandler(
            http.Dir("/static/"),
            "1234567" // find the commitish here
            nil,
            false,
    )

    http.ListenAndServe(":8080", &handler)

Documentation
-------------

- [Online Documentation (godoc.org)](http://godoc.org/github.com/ant0ine/go-static-forever)

Copyright (c) 2013 Antoine Imbert

[MIT License](https://github.com/ant0ine/go-static-forever/blob/master/LICENSE)


