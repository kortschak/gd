# `gd`

The `gd` program runs a single file Go program and renders it and its output as Markdown. It also allows inclusion of plain prose via C-style comments with a `{md}` prefix.

## Example output

This is the output from running the Go Hello World program with `gd`.

````
```
package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
```
> ```stdout
> Hello, world!
> ```
```
}
```
````

Which will render as this:

```
package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
```
> ```stdout
> Hello, world!
> ```
```
}
```

`gd` can also [include graphic output](examples/images) in the Markdown document.

## Limitations

Panic calls cannot be reflected in output since a panic takes control away from the `gd`-running program. To be able to capture output and associate it with source code lines, `gd` rewrites imports of "fmt" and "log" to "github.com/kortschak/gd/fmt" and "github.com/kortschak/gd/log". Behaviour of "fmt" is well replicated, but `gd` replaces `log.Panic*` calls with a simulation of a panic that outputs a stack trace and then exits. This means that stack unwinding is not performed and a `log.Panic*` call cannot be recovered. The `panic` built-in behaves as normal, but the panic output cannot be retrieved by `gd`.