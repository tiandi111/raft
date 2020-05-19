Lesson Learned
---

#### 1. Abstract state machine transition

To me, state transition is the most important aspect in raft algorithm. One 
lesson I've learned from this project is to abstract the state transitions
to independent, closed methods on go struct. Here are reasons I do this:

① It forces you to picture the whole lifecycle of a raft node as a state
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
when data inconsistency is detected and so on.

② It simplifies your code and helps you avoid bugs. Raft state includes a 
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

③ Improved readability. In ② I recommend you to encapsulate state transitions
to independent methods. You can also name each method to explicitly indicate the 
transition, PromoteToLeader, DemoteToFollower e.g. This also improves your code
readability.
 
#### 2. Think Term as logic clock

#### 3. Restrict lock(and channel) usage in as fewer sections as possible

#### 4. Client-Server separation