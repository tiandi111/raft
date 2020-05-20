Lessons Learned
---

#### 1. Abstract state machine transition.

To me, state transition is the most important aspect in raft algorithm. One 
lesson I've learned from this project is to abstract the state transitions
to independent, closed methods and attached on  go struct. Here are reasons 
I do this:

- It forces you to picture the whole lifecycle of a raft node as a state
machine. As a beginner, I feel hopeless when deal with distributed things 
because of its complexity, unpredictability, indefiniteness and asynchronous 
nature. However with a state machine, you clearly know all the possibilities
(with no fault happens) that your raft node should go in the next step so to 
unload your mind. Once you wrap up your algorithm under ideal state, you
can take a step further to think about how to tolerant fault, handle data
inconsistency and other problems you might meet in distributed environment.
From my experience, once you work out the state machine model and encapsulate
all state transitions to independent methods, you can easily solve those
distributed problems by making the transition atomic, abort certain actions
when data inconsistency detected and so on.

- It simplifies your code and helps you avoid bugs. Raft state includes a 
dozen of fields and they all closely related. For example, term is the core 
state in raft since it defines the logic clock of the whole system, theoretically
every other state variables should have a corresponding term with it. If you
enter a new term and you forget to update some state variables your cluster
might behave unpredictably. For example, if for some reasons your cluster 
didn't choose a new leader in a term and you have two candidates. If you forget
to update the number of votes they received then in the next election phase
the two candidate will both promote themselves to leader which is undesirable.
To avoid this, you can encapsulate every possible state transition to an independent
method and expose only the fields that you need to update corresponding to this
transition and hide the fields that should be assigned with default or zero value. 

- Improved readability. Last section I recommend you to encapsulate state transitions
to independent methods. You can also name each method to explicitly indicate the 
transition, PromoteToLeader, DemoteToFollower e.g. This also improves your code
readability.
 
#### 2. Think Term as logic clock.

Term is the timestamp in raft cluster. If you think the state of the cluster as 
an append-only log, then each log should have a timestamp attached to it, which 
is term in raft.

#### 3. Restrict lock(and channel) usage in as fewer sections as possible.

I've learned raft algorithm the hard way, use locks and channels everywhere
in my code, at the beginning. This makes thing difficult in three ways. 

- First, your code is hard to read because of those lock() calls and unlock() calls 
that scattered everywhere. 

- Second, dead-lock is very likely to happen but very unlikely to find especially
when you lock your code and call another function which acquires the same lock.

- Third, resource leak is likely to happen. 

The way I solved the above problems is to restrict lock and channel usage in only a
few sections. By coding the algorithm, I find out that most lock usages happen
when the node is going to take some activities, and it needs to ensure that it is 
still the same state as the time it made the decision to take those activities. Hence, 
I extract this part of logic to a method EnsureAndDo(). It locks the node status, decide 
if the passed in functions should be executed, unlock the status and return the result
as a bool.
```go
func (n *Node) EnsureAndDo(term, state int32, fs ...func(node *Node)) bool {
	n.LockStatus()
	defer n.UnlockStatus()

	if n.CurrentTerm == term && n.State == state {
		for _, f := range fs {
			f(n)
		}
		return true
	}

	return false
}
```
The only thing to note is that the functions passed in should not call EnsureAndDo() 
inside and this is easy to guarantee. By doing this, I removed almost all lock() and 
unlock() call in my code. Also, when dead-lock do happens, it is easy to debug since 
you know there are only a few places that calls EnsureAndDo that might cause the dead-lock. 
Simpler and safer.

I only use channel to inform a node that it has been elected as leader. However, inserting
channel code to the main logic causes the same code to be repeated multiple times. Hence,
I restrict all channel code into state transition functions. This again makes my code
simpler and safer.

#### 4. Client-Server separation

In raft, each node serve as both client and server to others. Mixing client 
functionalities and server functionalities makes it hard to code and read. 
I use grpc, a rpc framework, that naturally separates your services into client 
part and server part. In addition, golang's good support to multi-concurrency
programming and inter-goroutine communication(channel) makes it easier to 
decouple client and server code.
