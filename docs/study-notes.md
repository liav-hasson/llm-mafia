First document: https://go.dev/doc/effective_go#concurrency

> [!NOTE]
> Draft document!

#### Channels
- Channels in go are pipes that allow goroutines to send and receive values to each other. this method is the way concurrency is implemented in Go.
- a channel is always typed - meaning it has a datatype (chan int, chan string, etc..)
- send and receive operations are blocking, meaning that senders and receivers wait on the channel until their counterpart is ready.
- Channel have 2 main types:
unbuffered channels: 0 capacity, default, force synchronous communication. think - the meeting point of 2 goroutines.
this type means that both goroutines' execution positions are known at the value exchange point.
buffered channels: set capacity (e.g. 'make(chan int, 10)')
#### Channel syntax
- must use make() to create:
ch := make(chan int)
- Use the arrow operator to send data to the channel:
send operation: ch <- 42
- Set a variable by receving a value fron the channel:
v := <-ch
- The Close() method:
closes the channel, meaning no more variables can be sent.
also, it keeps the buffer to allow receivers to read values.
#### Goroutines
- Goroutines and channels are closely related.
- they use the same model of the unix pipe: producer | transformer | consumer. go's model originates in Hoare's (old commputer science guy) Communicating Sequential Processes (CSP).
- goroutines are different then threads, coroutines, processes.
- goroutines are multiplexed to multiple os threads, so if one is blocking a thread, another can continue to run.
- goroutines are more lightweight then threads, costing little more than the allocation of the stack space. (the stack starts small and can grow).
- goroutines run in the same address space as others.
- to run a function with a goroutine - prefix it with 'go' (e.g. 'go func()')
