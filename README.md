# chans

chans提供了一组常用chan操作：

- `Closed`：判断chan是否已关闭；
- `CloseOnce`：确保chan只被关闭一次，可多次调用`CloseOnce`，但chan只被关闭一次；
- `Done`：常用于`if Done(ctx.Done()) { ... }`模式。

另外，chans还提供了一种类似内置chan，但不限容量的`Chan[T]`结构，所有行为（包括阻塞和panic行为）对标内置chan：

- `New[T]()`：类似`make(chan T)`，没有capacity参数，但没有容量限制；
- `Send`：类似`ch <- value`，但因没有容量限制，写入不会阻塞；
- `Recv`：类似`value := <-ch`；
- `TryRecv`：类似`value, ok := <-ch`；
- `Close`：类似`close(ch)`。



## 示例

```go
package main

import (
	"context"
	"fmt"

	"github.com/eachain/chans"
)

func builtin() {
	var ch chan int
	fmt.Println(chans.Closed(ch)) // Output: false
	ch = make(chan int)
	fmt.Println(chans.Closed(ch)) // Output: false
	chans.CloseOnce(ch)
	fmt.Println(chans.Closed(ch)) // Output: true

	chans.CloseOnce(ch) // close multi times, will not panic

	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println(chans.Done(ctx.Done())) // Output: false
	cancel()
	fmt.Println(chans.Done(ctx.Done())) // Output: true
}

func mychan() {
	ch := chans.New[int]()
	ch.Send(123)
	fmt.Println(ch.Recv()) // Output: 123
	ch.Send(456)
	ch.Close()
	fmt.Println(ch.Recv()) // Output: 456
	fmt.Println(ch.Recv()) // Output: 0
}

func main() {
	builtin()
	mychan()
}
```

