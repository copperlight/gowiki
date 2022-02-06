## Introduction

Implementation of other tasks from the Go [Writing Web Applications](https://go.dev/doc/articles/wiki/)
Tutorial:

* Store templates in `tmpl/` and page data in `data/`.
* Add a handler to make the web root redirect to `/view/FrontPage`.
* Spruce up the page templates by making them valid HTML and adding some CSS rules.
* Implement inter-page linking by converting instances of `[PageName]` to
`<a href="/view/PageName">PageName</a>`. (hint: you could use `regexp.ReplaceAllFunc` to do this).

## Local Development

```bash
go build wiki.go
./wiki
open http://localhost:8080
```
