Mutual Exclusion Object or Mutex for short.   
  
To understand what mutex is lets take an example.  
lets assume that while we are running programs lets say that there are 2 threads "THREAD A" and "THREAD B"  
these threads are performing very similar tasks in which they have take the value of a common variable,   
Let's call this as count as of now count value is 1000 and both of the threads have to decrease counts value.   
So lets say that both threads start at the same time.  
Now THREAD-A takes the value of the count and then Decrease it to 999 and then THREAD-B takes this value and then decrease it to 998,   
but there what could happen is that THREAD-B instead of Taking value 999 it takes 1000(i.e. before THREAD-A finishes the task), 
and this can cause a lot of problems in our program, So to counter this we use Mutexes.  
A Mutex is like a lock, a THREAD applies it to some of the variables while it is doing a task with the help of mutex and   
then it prevent any other THREAD to access the same locked variables.  
Now here comes the problem. What if 2 or more threads try to apply the lock at the same time?   
So for this, a good application of mutex counters this by providing 2 states one locked and unlocked,   
now we check in a loop which constantly runs while the task is not completed. It checks again and   
again whether the mutex is locked or unlocked, if it is locked then nothing happens and the loops runs again,   
if it is unlocked it locks it and then allows the thread to perform the task. 
Now for golang, there are Goruotine instead of threads.
Let's say that there are some public variables that all goroutines can access,   
and our program uses these public variables for some task so one of our goroutine locks the variables and   
then performs its task such that the other goroutine access the variable after the first one has finished the task.
--
How to use  
--
I wont go into much details into code because i dont understand it that much. but implementation of it seems easy   
we just need to import fmt and sync pakages  
and then assign a var to sync.Mutex  
then use that assigned variable with .lock() to lock it and then perform our task and once our task is complete we do .unlock() to unlock it   
