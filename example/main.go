package main

import(
        "github.com/ant0ine/go-static-forever"
        "net/http"
)

func main() {
        handler := forever.NewStaticHandler(
                http.Dir("."),
                "1234567",
                nil,
                false,
        )
        http.ListenAndServe(":8080", handler)
}
