# go-caching

This is a project that replicates some of the features of the redis and memcached caching system, written in Go. It works with consistent hashing and can run on multiple servers (as well as responding to when a server is down). Currently a work in progress.

To run it, simply run `server_node.sh` and pass in the desired port. Then, run `client.sh` and pass in the ports of the server nodes. The file `test.sh` runs five server nodes on one machine.
