Test plan
---

Some cases
---
1. Run the cluster in a fault-free environment 
2. Run the cluster in a fault-free environment, then manually crash the leader
3. Run the cluster in a fault-free environment, then manually crash a follower
    - todo: all node should remove crashed nodes
4. Run the cluster in a fault-free environment, then ignore AppendEntriesRequest to mock network problems 
    - see if the follower can catch up its logs with the leader 
 
Fault injection
---
How?