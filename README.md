# Ignore

A simple ignore file processor for Go.

This is a reimagining of [jimschubert/iggy](https://github.com/jimschubert/iggy).

## Usage

Given an ignore file, for example `/your/directory/.gitignore`:

```
# This is an example ignore file

# Match everything below this directory
path/to/dir/**

# Match all Test files with all extensions
**/*Test*.*

# Match files with one character filename and extension
**/?.?

# Match all files beginning with first or second
**/{first*,second*}
```

This file can be read and evaluated like this:

```go
processor, _ := NewProcessor()
var allow bool
allow, _ = processor.AllowsFile("path/to/dir/ignored") // = false
allow, _ = processor.AllowsFile("other/path/to/dir/ignored") // = false
allow, _ = processor.AllowsFile("nested/test/SomeTest.java") // = false
allow, _ = processor.AllowsFile("nested/a.b") // = false
allow, _ = processor.AllowsFile("nested/abc.d") // = true
allow, _ = processor.AllowsFile("nested/first.txt") // = false
allow, _ = processor.AllowsFile("nested/second.txt") // = false
allow, _ = processor.AllowsFile("nested/third.txt") // = true
```

The `NewProcessor` function accepts functional parameters as defined in the ignore package.

To specify a custom path for the ignore file:

```go
processor, _ := NewProcessor(
    WithIgnoreFilePath("/your/directory/.ignore"),
)
```

Currently only Git's ignore rules are supported. As new strategies are added they can be targeted like this:

```go
processor, _ := NewProcessor(
    WithGitignoreStrategy(),
    WithIgnoreFilePath("/your/directory/.ignore"),
)
```

## Patterns

File patterns of the default ignore strategy follow closely to that of `.gitignore`.

* Rooted file pattern: `/*.ext`
    - Must exist in the root of the directory
    - Must begin with a forward slash `/`
    - Supports `*` or `*.ext` pattern
* Directory Rule
    - Matches against directories (`dir/`) or directory contents (`dir/**`)
    - Must end in `/`
* File Rule
    - Matches an individual `filename` or `filename.ext`

Similar to `.gitignore` processing, a double asterisk (`**`) can be used in place of a directory to indicate recursion.

For example:

```
path\to\**\file
```

matches both `path\to\some\file` and `path\to\some\nested\file`.

Single asterisks (`*`) match any characters within a pattern.

For example:

```
path\to\*file
```

matches both `path\to\your_file` and `path\to\my_file`, as well as `path\to\file`.

## Why?

I mean… why not? Sometimes I want a simple way to ignore or force file processing in a directory, but I don't want to shell out to some other program to evaluate the logic.

# License

Apache 2.0.

see [License](./LICENSE)
