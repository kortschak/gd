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
