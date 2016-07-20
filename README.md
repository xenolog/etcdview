# EtcdView
Utility for comfortable operation with data into etcd cluster

## Usage
This is a command-line utility, which can make beter you operation with data into etcd cluster. Use `etcdview help` for getting help about commands and options.

* display `tree` (or subtree) of data from existing etcd cluster

        $ ./etcdview -u http://127.0.0.1:4001 tree /calico/bgp
        calico:
          bgp:
            v1:
              global:
                as_num: 64511
                node_mesh: {"enabled": true}
              host:
                node-5.domain.local:
                  ip_addr_v4: 10.88.11.7
                  ip_addr_v6:
                node-6.domain.local:
                  ip_addr_v4: 10.88.11.5
                  ip_addr_v6:
                node-7.domain.local:
                  ip_addr_v4: 10.88.11.11
                  ip_addr_v6:
        $
     
## Build
Go 1.5+ required

* Setup Go development environment (like http://skife.org/golang/2013/03/24/go_dev_env.html for example)
* cd $GOPATH
* go get github.com/xenolog/etcdview
* cd src/github.com/xenolog/etcdview/
* go get
* go build

Static linked binary `etcdview` should be present into current dir.
