# cache-inventory
Inventory cache for qframe, to enable queries against an in-memory inventory snapshot.

## Integration test

The script within `cmd/main.go` uses a request to query for a given container name and reuses the respons' IP to query the same container by his IP.

### Container using a bridge

The container was started like this: `docker run -ti --rm --name bridge-cnt ubuntu tail -f /dev/null`

```bash
$ go run main.go bridge-cnt
[  INFO] Dispatch broadcast for Data, Done and Tick
[NOTICE]       inventory Name:inventory  >> Start inventory v0.3.1
[  INFO]       inventory Name:inventory  >> Create query for container by Name 'bridge-cnt' : qcache_inventory.NewNameContainerRequest('q1', 'bridge-cnt')
[NOTICE]   docker-events Name:docker-events >> Start docker-events collector v0.3.0
[ DEBUG]       inventory Name:inventory  >> Received InventoryRequest for {2017-08-29 15:21:56.691092464 +0000 UTC m=+1.009847856 q1 1s bridge-cnt   0xc420204840}
[  INFO]   docker-events Name:docker-events >> Connected to 'moby' / v'17.07.0-ce-rc3'
[ DEBUG]       inventory Name:inventory  >> Add CntID:ad4bbcdfd71c5 into Inventory (name:/bridge-cnt, IPs:172.17.0.3)
[ DEBUG]       inventory Name:inventory  >> Add CntID:4a4b8d12da535 into Inventory (name:/modest_einstein, IPs:172.17.0.2)
[  INFO]       inventory Name:inventory  >> Got InventoryResponse: Container 'bridge-cnt' has ID-digest:ad4bbcdfd71c5
[  INFO]       inventory Name:inventory  >> Use IP from first query response container (network:bridge) to generate another query: qcache_inventory.NewIPContainerRequest('q2', '172.17.0.3')
[ DEBUG]       inventory Name:inventory  >> Received InventoryRequest for {2017-08-29 15:21:57.190813115 +0000 UTC m=+1.509569306 q2 2s   172.17.0.3 0xc4202c2900}
[  INFO]       inventory Name:inventory  >> Got InventoryResponse: Container w/ IP 172.17.0.3 has Digest:ad4bbcdfd71c5
```

### Container connected to a network.

Container start:

```bash
$ docker network create testnet --subnet=192.168.0.0/16
b804cd1f408c7c757b2541d5ccaf75991c18d7ec005be608b774c9b49929f9a4
$ docker run -ti --rm --name testnet-cnt --network testnet ubuntu bash
root@eb1da0b083f0:/#
```

Run the tool.

```bash
$ go run main.go testnet-cnt
[  INFO] Dispatch broadcast for Data, Done and Tick
[NOTICE]       inventory Name:inventory  >> Start inventory v0.3.1
[  INFO]       inventory Name:inventory  >> Create query for container by Name 'testnet-cnt' : qcache_inventory.NewNameContainerRequest('q1', 'testnet-cnt')
[ DEBUG]       inventory Name:inventory  >> Received InventoryRequest for {2017-08-29 15:26:57.896243684 +0000 UTC m=+1.011992256 q1 1s testnet-cnt   0xc420204a80}
[NOTICE]   docker-events Name:docker-events >> Start docker-events collector v0.3.0
[  INFO]   docker-events Name:docker-events >> Connected to 'moby' / v'17.07.0-ce-rc3'
[ DEBUG]       inventory Name:inventory  >> Add CntID:eb1da0b083f01 into Inventory (name:/testnet-cnt, IPs:192.168.0.2)
[ DEBUG]       inventory Name:inventory  >> Add CntID:4a4b8d12da535 into Inventory (name:/modest_einstein, IPs:172.17.0.2)
[  INFO]       inventory Name:inventory  >> Got InventoryResponse: Container 'testnet-cnt' has ID-digest:eb1da0b083f01
[  INFO]       inventory Name:inventory  >> Use IP from first query response container (network:testnet) to generate another query: qcache_inventory.NewIPContainerRequest('q2', '192.168.0.2')
[ DEBUG]       inventory Name:inventory  >> Received InventoryRequest for {2017-08-29 15:26:58.404090241 +0000 UTC m=+1.519842310 q2 2s   192.168.0.2 0xc4202c27e0}
[  INFO]       inventory Name:inventory  >> Got InventoryResponse: Container w/ IP 192.168.0.2 has Digest:eb1da0b083f01
```
