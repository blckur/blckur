[![blckur](media/logo.png)](https://github.com/blckur/blckur)

# blckur: in development

Copyright (c) 2015 Zachary Huff

[![godoc](https://godoc.org/github.com/blckur/blckur?status.png)](https://godoc.org/github.com/blckur/blckur)

## Node types

* app - web server
* cache - redis server
* queue - beanstalkd server
* scheduler - send scheduled tasks to queue
* worker - run queue jobs

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
