# kafka-container-source

A simple, pluggable, kafka-based event source to be used in a Knative Eventing context.

## Usage

To create a source, reading from a given Apache Kafka cluster, do a `ko apply -f` with a YAML like the following:

```yaml
apiVersion: sources.eventing.knative.dev/v1alpha1
kind: ContainerSource
metadata:
  name: my-kafka-source
spec:
  image: github.com/matzew/kafka-container-source/cmd
  args: 
    - '--bootstrap=kafkabroker.kafka:9092'
    - '--topic=my-topic'
    - '--groupId=foobar'
  sink:
    apiVersion: eventing.knative.dev/v1alpha1
    kind: Channel
    name: testchannel
```


### Debugging

#### local build ...

```
go build -o kes ./cmd/main.go
```

and execute it:

```
./kes --bootstrap=localhost:9092 --topic=my-topic --groupId=test25 --sink=http://localhost:7080

```

#### Docker build

```
docker build --no-cache -t docker.io/username/knative-kafka-src:latest .
```

