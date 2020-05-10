Design Draft
---

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
func handle_request_vote():
	while:
	    if vote_request = recv():
	        if vote_request.term < term || (has_voted && vote_request.term<vote_term):
	            continue
	        else:
	            server = vote_request.server
	            vote(server)
	            has_voted = true
	            vote_term = vote_request.term
	            
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

