# sarama-rack

A test rig for https://github.com/Shopify/sarama/pull/1696.

You'll need docker-compose and all that jazz.

Edit docker-compose.yml and set `KAFKA_ADVERTISED_HOST_NAME` to the address returned by `ifconfig en0` or so on a Mac.

Then `docker-compose up --scale kafka=2 -d`. It'll take 30 seconds or so for things to be ready. Check `docker-compose logs`.

Then run the test app:

```
go run main.go 192.168.1.10:$(docker port sarama-rack_kafka_1 9092 | cut -d: -f2)
```

Replace `192.168.1.10` with what you set `KAFKA_ADVERTISED_HOST_NAME` to.

When the test app starts up, it finds the rack name of the non-leader replica. It then connects to consume from the cluster with that rack name.

It may take a few runs (once the replica is in sync), but eventually you should see something like:

```
...
2020/10/04 18:26:18 replica ID 1002 rack 2f6459093c13
...
[sarama] 2020/10/04 18:26:18 producer/broker/1001 starting up
[sarama] 2020/10/04 18:26:18 producer/broker/1001 state change to [open] on test-1/0
[sarama] 2020/10/04 18:26:18 Connected to broker at 192.168.1.10:32826 (registered as #1001)
[sarama] 2020/10/04 18:26:18 consumer/broker/1001 added subscription to test-1/0
[sarama] 2020/10/04 18:26:21 consumer/test-1/0 finding new broker
[sarama] 2020/10/04 18:26:21 client/metadata fetching metadata for [test-1] from broker 192.168.1.10:32826
[sarama] 2020/10/04 18:26:21 ClientID is the default of 'sarama', you should consider setting it to something application-specific.
[sarama] 2020/10/04 18:26:21 consumer/broker/1002 added subscription to test-1/0
```

Here, the initial connection was to the leader (1001). It returned a PreferredReadReplica of the replica (1002) and the client reconnected to that.
