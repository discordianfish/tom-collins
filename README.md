# tom - our collins cli
*3 parts Old Tom Gin, 2 parts lemon juice, 1 part sugar syrup, 4 parts carbonated water to taste*

## Usage
```
tom v0.0.0-dev - Usage: tom [options] sub-command
  -password string
      Collins password (default "admin:first")
  -uri string
      URL to Collins API (default "http://localhost:9000/api")
  -user string
      Collins user (default "blake")
  -v  Print version and exit

Sub commands:
  - query <cql query>
   query runs a CQL query and returns matching assets
  - register <hostname>
   register is to be implement
  - tag <hostname>
   tag returns tag for hostname
  - template <template:[destination]> [template...]
   template renders templates with assets returned by given query.
  - update <attribute> [value]
   update Update attributes on assets
```

### Sub commands
#### Template
The templates specified can be either local files or collins attributes.
In the template `.Assets` contains a list of all Assets in collins which were
returned for the query (See `-q`), including all their attributes. See
[go-collins documentation](https://godoc.org/github.com/tumblr/go-collins/collins#Asset) for
available attributes.

`tom template foo.txt.tmpl:foo.txt` renders the local template foo.txt.tmpl and
writes the result to foo.txt. If the source template has a suffix `.tmpl`, the
destination can be ommited and the output will be written to the full path of
the template with the `.tmpl` suffix stripped. Therefor this will create the
same file as above: `tom template foo.txt.tmpl`.

`tom template -r foo:foo.txt` would use the attribute `foo` from asset `default`
as template source. If the source contains a slash, the string before the slash
is used to select a different asset: `tom template -r bar/foo:foo.txt` uses the
attribute `foo` from asset `bar` and writes to foo.txt.
