## Installation

    $ go get github.com/mycalf/reader-go


## Examples

```Go
package main

import (
	"fmt"

	"github.com/mycalf/reader-go"
)

func main() {
	if doc, ok := reader.Load("http://www.163.com"); ok {
		fmt.Println(doc.HTML)
	}
}
```