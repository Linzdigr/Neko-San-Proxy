package main

import (
    "io"
    "net/http"
    "net/url"
)

func proxy(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    urlStr := r.Form.Get("u")

    urlParsed := url.Parse(urlStr)

    w.Header().Set("X-Proxy-Go", "Neko-San-Proxy-v1")

    if r.Method == "GET" {
        response, err := http.Get(urlStr)
        if err != nil {
            errorOut(err)
        } else {
            defer response.Body.Close()
            contents, err := ioutil.ReadAll(response.Body)
            if err != nil {
                errorOut(err)
            }
            fmt.Fprintf(w, "%s\n", contents)
        }
    }
}

func main() {
    http.HandleFunc("/", proxy)
    http.ListenAndServe(":8008", nil)
}

/*
    type URL struct {
            Scheme     string
            Opaque     string    // encoded opaque data
            User       *Userinfo // username and password information
            Host       string    // host or host:port
            Path       string
            RawPath    string // encoded path hint (Go 1.5 and later only; see EscapedPath method)
            ForceQuery bool   // append a query ('?') even if RawQuery is empty
            RawQuery   string // encoded query values, without '?'
            Fragment   string // fragment for references, without '#'
    }
*/
