# jsonexec

Run command and convert output (stdout) JSON into provided object

## How to use

```go
package main

import (
	"fmt"
	"github.com/sirkon/jsonexec"
)

func main() {
	var dest map[string]interface{}
	if err := jsonexec.Run(&dest, "ls"); err != nil {
		jsonexec.HandleError(err, func(lsOutput string) {
			// will show something like
			// unmarshal command output: invalid character 'L' looking for beginning of value
			// Command output:
			// LICENCE
			// README.md
			// ...
			fmt.Println(err, "\nCommand output:\n"+lsOutput)
		})
	}
}
```
