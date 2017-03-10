# Description

This is a light proxy which checks locally cluster (mariadb, percona, mysql, rabbitmq) and returns http code based on current status:
- 200 - Everything is fine. Node in cluster
- 503 - Node is out of sync

Intention of this service is to check consistency of a cluster.  
If node gets another state but synced, this daemon will return 5xx to a http request from haproxy, keepalived or another service.  
This will let you to exclude node from LBpool.   

# Usage

- Edit configuration file /etc/clusterhc/clusterhc.toml  
- Run service
- Connect configure your LBpool to check *localBind*/(galera|rabbimq)

# Installation

- Install go https://golang.org/doc/install
- Make a proper structure of directories: ```mkdir -p /opt/go/src /opt/go/bin /opt/go/pkg```
- Setup g GOPATH variable: ```export GOPATH=/opt/go```
- Clone this project to src: ```go get github.com/leoleovich/clusterhc```
- Fetch dependencies: ```cd /opt/go/github.com/leoleovich/clusterhc && go get ./...```
- Compile project: ```go install github.com/leoleovich/clusterhc```
- Copy config file: ```mkdir /etc/clusterhc && cp /opt/go/src/github.com/leoleovich/clusterhc/clusterhc.toml /etc/clusterhc/```
- Run it ```/opt/go/bin/clusterhc```

# Optional systemd service
- Create systemd service: ```cp /opt/go/src/github.com/leoleovich/clusterhc/clusterhc.service /usr/lib/systemd/system```
- Reload systemd: ```systemctl daemon-reload```
- Start clusterhc service: ```systemctl start clusterhc.service``` 
- Check status: ```systemctl status clusterhc.service```