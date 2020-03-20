# pubsub2syslog

A simple go-language program that can read from a Pub/Sub topic and write to syslog (udp/tcp).

## Usage

```
pubsub2syslog -project your-project -server syslog-server:514 -protocol tcp -topic pubsub-topic
```

