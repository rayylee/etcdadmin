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
# go get -u github.com/rayylee/etcdadmin/etcdadmind
# go get -u github.com/rayylee/etcdadmin/etcdadminctl
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
