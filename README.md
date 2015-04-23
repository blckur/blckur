[![blckur](media/logo.png)](https://github.com/blckur/blckur)

# blckur: in development

Copyright (c) 2015 Zachary Huff

[![godoc](https://godoc.org/github.com/blckur/blckur?status.png)](https://godoc.org/github.com/blckur/blckur)

## Node types

**type**    | app
----------- | -----------
**use**     | gin web server
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/app/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/app)

**type**    | cache
----------- | -----------
**use**     | redis server
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/cache/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/cache)

**type**    | queue
----------- | -----------
**use**     | beanstalkd server
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/queue/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/queue)

**type**    | queue
----------- | -----------
**use**     | beanstalkd server
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/queue/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/queue)

**type**    | scheduler
----------- | -----------
**use**     | schedule tasks in queue
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/scheduler/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/scheduler)

**type**    | worker
----------- | -----------
**use**     | run tasks in queue
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/worker/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/worker)

## Systems

### cache

Redis cluster with client side sharding/replication

* fast
* low consistency
* high availability
* partition tolerant
* cache only no guarantees
* supports get/set/pub/sub

### queue

Beanstalkd cluster with client side sharding/replication

* fast
* low consistency
* high availability
* partition tolerant
* send queue jobs to worker nodes

### messenger

Pub/sub messaging system using mongodb tailable cursor

* slow
* consistent
* all nodes read messages
* send system messages/cluster change notifications to nodes
