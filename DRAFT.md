Design Draft
---

## RPCs
- AppendEntries RPC
```
func (term, success)  AppendEntries(term, leaderID, prevLogIndex, 
    prevLogTerm, entries[], leaderCommit):
    
    // 1. stale call
    if term < currentTerm:
        return (currentTerm, false)    
        
    // 2. lack of previous logs, cannot apply
    if localLogs.len <= prevLogIndex:
        return (currentTerm, false)
        
    // 3. find inconsistent entry, delete entries followed
    int startFrom
    for i from prevLogIndex+1 to min(localLogs.len-1, entries.len):
        if localLogs[i].term != entries[i-prevLogIndex-1].term:
            startFrom = i-prevLogIndex-1
            delete(localLogs, i, localLogs.len)
            break
          
    // 4. append any new entries not already in the log  
    for i in (startFrom, localLogs.len):
        append(localLogs, entries[i])
        
    // 5. set commitIndex
    // Here we set the minimum of leaderCommit and localLogs.len-1 
    // because we don't know yet if this AppendEntris call will 
    // succeed in the mojority of the nodes
    if laderCommit > commitIndex:
        commitIndex = min(leaderCommit, localLogs.len-1)    

```

## Part1: Leader Electrion
### Leader
- Heartbeat
```
func heartbeat() {
	while:
	    for each follower:
	    	send(heartbeat_signal)
	    sleep(heartbeat_interval)
}
```

### Follower
- Heartbeat Daemon
```
func heartbeat_daemon():
    while:
    	if heartbeat_signal = recv(); hearbeat_signal != null :
    	    if heartbeat_signal.term >= term:
                last_heartbeat = heartbeat_signal
                last_heartbeat.time = time.now()
                term = heartbeat_signal.term
```
- Heartbeat Monitor
```
func heartbeat_monitor():
	while:
	    if time.now() - last_heartbeat.time > election_timeout:
	        state = CANDIDATE
	    	begin_election()
	    	return
	    sleep(minitor_interval)
```
- Request Vote Handler
```
func boolean handle_request_vote():
	while:
	    if vote_request = recv():
	        if vote_request.term < term || (has_voted && vote_request.term >= term):
	            return false
	        else:
	            server = vote_request.server
	            has_voted = true
	            vote_term = vote_request.term
	            return true
```

### CANDIDATE
- Begin Election
```
func boolean begin_election():
    while:
        term = term+1
        vote(self)
        start = time.now
        received_votes = 0
        for each server in cluster:
            async_request_vote(server, term, received_votes)
        while time.now() - start < election_timeout:
            if received_votes > total_num_nodes/2:
                state = LEADER
                heartbeat()
                return true
            else if heartbeat_signal = recv(); heartbeat_signal != null :
                if heartbeat_signal.term >= term:
                    state = FOLLOWER
                    return true
                else:
                    ignore(heartbeat_signal)
                    continue
            else:
                // no result, do nothing
```
- Request Vote
```
func async_request_vote(server, term, received_votes):
	new thread execute:
	    resp = request_vote(server, term)
	    if resp == GOT_VOTE:
	    	received_votes += 1 // atomic
	return 
```
## Part2: Log Replication

## Part3: Safety Mechanism

