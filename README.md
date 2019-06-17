# list-imports

`list-imports` is a Go command line tool for listing the imports found in a Go project.

```
List the imports of a go project or specified directory as JSON. Optionally shows
imports of vendor folder, essentially showing transitive dependencies.
Specifying a dir is optional, and will default to the working directory.

Usage:
  list-imports [dir] [flags]

Flags:
  -h, --help             help for list-imports
  -i, --include-vendor   include vendor dir in listing of imports
```