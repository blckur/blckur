[![blckur](media/logo.png)](https://github.com/blckur/blckur)

# blckur: simple notifications

Copyright (c) 2015 Zachary Huff

[![godoc](https://godoc.org/github.com/blckur/blckur?status.png)](https://godoc.org/github.com/blckur/blckur)

## node types

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

**type**    | scheduler
----------- | -----------
**use**     | schedule tasks in queue
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/scheduler/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/scheduler)

**type**    | worker
----------- | -----------
**use**     | run tasks in queue
**docker**  | [![Docker Repository on Quay.io](https://quay.io/repository/blckur/worker/status "Docker Repository on Quay.io")](https://quay.io/repository/blckur/worker)

## systems

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

## aws nodes

**name** | restless-thunder-4475
-------- | --------
**use**  | elb load balancer
**dns**  | app.blckur.com

**name** | winter-waters-7121
-------- | --------
**use**  | mongodb replica node
**dns**  | winter-waters-7121.blckur.vpc

**name** | restless-forest-4412
-------- | --------
**use**  | mongodb replica node
**dns**  | restless-forest-4412.blckur.vpc

**name** | autumn-plains-4942
-------- | --------
**use**  | mongodb replica node
**dns**  | autumn-plains-4942.blckur.vpc

**name** | lively-forest-2840
-------- | --------
**use**  | mongodb replica set name

**name** | autumn-stars-6522
-------- | --------
**use**  | coreos fleet node
**dns**  | autumn-stars-6522.blckur.vpc

**name** | guarded-skies-3838
-------- | --------
**use**  | coreos fleet node
**dns**  | guarded-skies-3838.blckur.vpc

**name** | calm-plateau-4739
-------- | --------
**use**  | coreos fleet node
**dns**  | calm-plateau-4739.blckur.vpc

**name** | restless-waters-7160
-------- | --------
**use**  | coreos fleet node
**dns**  | restless-waters-7160.blckur.vpc

**name** | summer-skies-3064
-------- | --------
**use**  | coreos fleet node
**dns**  | summer-skies-3064.blckur.vpc

**name** | lively-waterfall-5258
-------- | --------
**use**  | coreos fleet node
**dns**  | lively-waterfall-5258.blckur.vpc
