# What is a Mutex ? (by b-iit : [https://github.com/b-iit])

## Introduction
Imagine you’re at a metro station, standing in a line trying to enter the metro, one by one. And then a guard is there which ensures a proper queue and allows everyone to enter inside the metro. In programming, a **mutex** (short for "mutual exclusion") works a lot like that guard. It ensures that only one thread (or "person") accesses a critical section (or "enter inside") at a time.

This is especially important in low-level programming languages like **Go (Golang)** or **Rust**, where you deal with **multi-threading** — allowing different parts of a program to run simultaneously. Multi-threading is everywhere when you need two or more operations to be performed simultaneously, i.e. in web servers handling multiple requests, databases managing transactions, or video games running physics and rendering together.

But, when multiple threads try to access and change the same piece of data at the same time, it leads to a problem called a **race condition**, which leads to subsequent errors in the execution of the program. Think of two persons trying to enter through a single door—someone might get hurt. That’s where mutexes and similar mechanisms come in to play.

---

## What is a Mutex?
A **mutex** is like a lock for your data. Here’s how it works:
1. When a thread wants to access a shared resource, it "locks" the mutex.
2. While the mutex is locked, no other thread can access the resource.
3. When the thread is done, it "unlocks" the mutex, allowing the next thread to use the resource.
   
Thus, it can be said that, mutual exclusion (mutex) is a property of concurrency control, which is implemented to prevent race conditions.

### Example in Everyday Life
- Imagine there is one Phone booth on a street:
- When someone enters, they close the door and talk on the phone.
- Others must wait until the call is done.
- This ensures privacy and order.

This explains how mutexes work.
In programming, mutexes are used to maintain this kind of order for shared resources. It becomes very important to enforce Mutex on uni-processor systems to ensure the program runs error-free.

---

## How to enforce Mutual Exclusion ?
Mutexes can be enforced in a program mainly two ways:
- ### Hardware solutions :
1) On uni-processor systems, the simplest solution to achieve mutual exclusion is to **disable interrupts** during a process's critical section. This will prevent any interruption to current service routines running.
2) **Busy-waiting technique** is effective for both uniprocessor and multiprocessor systems. In this technique shared memory and an atomic test-and-set instruction provide the mutual exclusion.
3) **Compare-And-Swap (CAS)** can be used to achieve wait-free mutual exclusion for any shared data structure by creating a linked list where each node represents the desired operation to be performed. CAS is then used to change the pointers in the linked list during the insertion of a new node. Only one process can be successful in its CAS; all other processes attempting to add a node at the same time will have to try again.

- ### Software solutions :
In addition to hardware-supported solutions, we can employ various algorithms to achieve mutual exclusion. Examples include:
1) Dekker's algorithm
2) Peterson's algorithm
3) Lamport's bakery algorithm
4) Szymański's algorithm
4) Taubenfeld's black-white bakery algorithm
5) Maekawa's algorithm

- It is often preferable to use synchronization facilities provided by an operating system's multithreading library, which will take advantage of hardware solutions if possible but will use software solutions if no hardware solutions exist. 


## Types of Mutexes

1) Semaphores
2) Readers–writer locks
3) Recursive locks
4) Monitors
5) Message passing
6) Tuple space

- Many forms of mutual exclusion have side-effects. For example, classic semaphores permit deadlocks, in which one process gets a semaphore, another process gets a second semaphore, and then both wait for the other semaphore to be released. Other common side-effects include starvation, in which a process never gets sufficient resources to run to completion.

---

## What is a Semaphore?
A **semaphore** is like a mutex but with more flexibility. While a mutex allows only one thread to access a resource, a semaphore can allow a specific number of threads to use the resource at the same time. It is a variable or abstract data type used to control access to a common resource by multiple threads and avoid critical section problems in a concurrent system such as a multitasking operating system. 
A semaphore is used to manage this kind of situation in programming. It tracks how many threads can access a resource simultaneously and allocates the resource to multiple threads based on capacity.

---

## Mutexes and Semaphores in Go (Golang)
Go provides built-in support for concurrency through **goroutines** (lightweight threads) and **channels** (communication between threads). To manage shared resources safely, you can use:

### Using Mutex in Go
The `sync` package provides a `Mutex` type. Here's an example:
```go
package main

import (
	"fmt"
	"sync"
)

var (
	counter int
	mutex   sync.Mutex
)

func increment() {
	mutex.Lock()
	counter++
	mutex.Unlock()
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

```

### Explanation
- **`counter int`**: Declares a variable counter of type int.
- **`mutex sync.Mutex`**: Declares a variable mutex of type sync.Mutex. sync.Mutex is a mutual extension lock in Go's sync package.
- **`mutex.Lock()`**: In the increment function, this locks the mutex. Now, other Goroutines can't lock it before it is unlocked.
- **`counter++ `**: Increases the value of counter variable by one when the increment function is called. 
- **`mutex.Unlock()`**: Unlocks the mutex after the increment is done.
- **`var wg sync.WaitGroup`**: Variable wg of type sync.WaitGroup to wait for a collection of Goroutines to finish.
- A for loop is used to create five Goroutines.
- **`wgAdd(1)`**: Increments WaitGroup counter to indicate new Goroutine is launched.
- **`defer wg.Done()`**: Signals the WaitGroup that this Goroutine's work is complete.
- **`wg.Wait()`**: Blocks the main function until all goroutines have signaled that they are done


- Without the mutex, multiple goroutines might read and write to `counter` simultaneously, causing incorrect results.

---

## Summary
- **Mutex**: Ensures that only one thread accesses a resource at a time. Use it for exclusive access.
- **Semaphore**: Allows multiple threads to access a resource while limiting how many can do so simultaneously. Use it for shared access with restrictions.
- **In Go**: Mutexes are provided by the `sync` package, and semaphores can be implemented using channels.

Understanding these tools will help you write safer and more efficient concurrent programs.
