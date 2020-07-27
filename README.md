# etcdadmin

etcdadmin contains `etcdadmind` and `etcdadminctl`. It makes it easy to create a new etcd cluster, add a member to, or remove a member from an existing etcd cluster.

- etcdadmind is a service to manage etcd cluster.
- etcdadminctl is a tool for operating etcdadmind, such as add, remove or list members.

## Depends

- etcd.service

Here is an example `contrib/etcd.service`

## Getting started

### How to install

```bash
# go install github.com/rayylee/etcdadmin/etcdadmind
# etcdadmind help
Usage: etcdadmind install | remove | start | stop | status

# go install github.com/rayylee/etcdadmin/etcdadminctl
# etcdadminctl help
A simple command line of etcdadminctl.

Usage:
  etcdadminctl [command]

Available Commands:
  help        Help about any command
  member      Membership related commands
  version     Prints the version

Flags:
      --endpoint string    gRPC endpoints (default "127.0.0.1:2390")
  -h, --help               help for etcdadminctl
      --write-out string   Output format (default "json")

```

### How to run
```bash
# etcdadmind install
Install Etcd admin service:                 [  OK  ]

# etcdadmind start
Starting Etcd admin service:                [  OK  ]

```
- Add member
    ```bash
    # etcdadminctl member add node1 --peer-ip=169.254.155.111
    # etcdadminctl member add node2 --peer-ip=169.254.155.112
    ```
- List memebers
    ```bash
    # etcdadminctl member list
    ```
- Remove member
    ```bash
    # etcdadminctl remove node1
    ```
