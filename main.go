// /!\ Lauch `go get golang.org/x/net/html` on a fresh installed server!

package main

import (
    "os"
    "io/ioutil"
    "fmt"
    "strings"
    "net/http"
    //"regexp"
    "golang.org/x/net/html"
    "net/url"
)

func errorOut(err error) {
    fmt.Printf("ERROR: %s", err)
    os.Exit(1)
}

func proxy(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    urlStr := r.Form.Get("u")

    urlParsed, err := url.Parse(urlStr)

    if err != nil {
        fmt.Printf("WARN: %s, Possibly due to Favico web browser request\n", err)
        return
    }

    w.Header().Set("X-Proxy", "Neko-San-Go-v1")

    if r.Method == "GET" {
        fmt.Printf("Requested  %s with scheme [%s]\n", urlParsed.Host, urlParsed.Scheme)
        response, err := http.Get(urlStr)
        if err != nil {
            fmt.Printf("WARN: %s, Possibly due to Favico web browser request\n", err)
            return
        } else {
            defer response.Body.Close()
            contents, err := ioutil.ReadAll(response.Body)
            if err != nil {
                errorOut(err)
            }

            DOMStr := string(contents[:])

            doc, err := html.Parse(strings.NewReader(DOMStr))
            if err != nil {
                fmt.Printf("WARN: %s\n", err)
            }

            output := DOMStr

            linkTab := []string{}

            var f func(*html.Node)
            f = func(n *html.Node) {
                if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "script" || n.Data == "link") {
                    for _, base := range n.Attr {
                        if base.Key == "href" {
                            linkTab = append(linkTab, base.Val)
                            fmt.Printf("Link A found: %s\n", base.Val)
                        } else if base.Key == "src" {
                            linkTab = append(linkTab, base.Val)
                            fmt.Printf("Link SRC found: %s\n", base.Val)
                        }
                    }
                }
                for c := n.FirstChild; c != nil; c = c.NextSibling {
                    f(c)
                }
            }

            f(doc)

            for _, el := range linkTab {
                if !strings.HasPrefix(el, "//") && !strings.HasPrefix(el, "http://") {
                    output = strings.Replace(output, el, "https://proxy.neko-san.fr/?u=http://" + urlParsed.Host + el, -1)
                    fmt.Printf("Replacement...\n")
                }
            }

            fmt.Fprintf(w, "%s", output)
            //html.Render(w, doc)
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
