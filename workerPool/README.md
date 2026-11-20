- Start 8 goroutines. Each one blocks on range jobChan, waiting for something to arrive.
- Fill the channel. The producer sends the IDs one at a time. As each one arrives, one of the idle workers takes it.
- Workers loop back. After finishing a job, a worker goes back to waiting on jobChan. That is why 8 workers can handle 200 jobs.

wg.Wait() is how the helper knows when it's safe to close resultChan.

What it waits for

Wait blocks until the WaitGroup counter reaches 0. The counter was raised to 8 by the wg.Add(1) calls, and each worker lowers it by one when it exits (defer wg.Done()). So Wait returns only after all 8 workers have finished and none of them can send again.

Why that's the right moment to close

Closing resultChan while a worker is still going to send on it causes a panic. After Wait returns, there are no senders left, so closing is safe. That's the whole job of this goroutine:

wait until all senders are done → then close the channel
