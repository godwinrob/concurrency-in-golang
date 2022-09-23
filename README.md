# Get GOing with concurrency in Go

## The Situation
 
We have two functions in need of a refactor to improve our web API.

For a solution we only want to use Go's standard library.

### Go Waitgroups
```func multipleDatabaseCalls()```

This function currently calls multiple databases one after the other.
The database calls do not depend on the returned results from the previous calls.

How can we utilize Go's Waitgroups and Go Routines to allow the database calls to be made simultaneously?

### Go Channels
```func functionWithHardTimeLimit()```

This function is running in AWS using Lambda functions with a hard timeout of 10 seconds.

Occasionally this function is taking longer than the 10 second limit, and is being killed by AWS.
When this happens all tracing is lost, so we are having trouble diagnosing the underlying cause.

How can we utilize Go's Channels and Go Routines to to ensure that it always returns before the Lambda timout so that the trace can be captured?