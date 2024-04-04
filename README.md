# LLM Feeder

This is a simple little tool to create an LLM-friendly file listing of a
project for questions like "what the fuck was the developer thinking when
building this garbage codebase??"

Totally not stolen from a friend and then enhanced to seem like a cool new
project of my own design.

## Ignoring stuff

The "-I" flag ignores by shell pattern, and it checks the relative path as well
as the name of an isolated directory entry. This is powerful but you have to
understand what it means.

A simple name with no special characters, like `-I .git`, will ignore only
things which have that exact name.

A longer path can be specified as well. If you want to ignore your vended
assets, for instance, you might specify `-I static/vendor`. That will never
match a filename unless you actually created a file with a literal slash in it,
so it's a more targeted approach than just ignoring "vendor", which may not be
desired.

All patterns are relative to the base directory. If you specify `-d foo/bar`,
and you want to exclude `foo/bar/static/vendor`, your ignore pattern is simply
`-I static/vendor`.

Shell patterns are powerful! You can ignore all XML files with a simple `-I
"*.xml"`. You can ignore XML that's in your test directory, but include XML
everywhere else, with something like `-I "tests/**/*.xml"`.

## Examples

Look at just the directory listing of a project (helpful when deciding on your ignore list)

```bash
./bin/feeder -d /path/to/something -q
```

Add a project, but skip various dirs/files you know you don't need analyzed

```bash
./bin/feeder -d /path/to/something -I ".git" -I "*.md" -I ".gitignore" -I "LICENSE"
```

Add only files git tracks while still ignoring certain undesired files

```bash
git ls-files | xargs ./bin/feeder -I "*.md" -I ".gitignore" -I "LICENSE"
```
