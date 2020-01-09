JSON Request Parser for the go (Golang)
===

Credit: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

This is a very small utility intended to properly parse JSON request and respond with various specific errors, if these occur.

Usage
---
```
import jparser "github.com/idoberko2/json-request-parser"

type request struct {
    Name string
    Age int
}

func (w http.ResponseWriter, r *http.Request) {
    var req request

	if ok := jparser.ParseJSONRequest(w, r, &req); !ok {
		return
	}

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Welcome, " + req.Name + "!"))
}
```
