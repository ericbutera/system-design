# Tilt notes

- [Search](#search)
  - [Opensearch](#opensearch)
- [NoSQL](#nosql)
  - [Cassandra](#cassandra)
- [RMDBS](#rmdbs)
  - [Postgres](#postgres)
- [Cache / Distributed Locks](#cache--distributed-locks)
  - [Redis](#redis)
- [Message Queue](#message-queue)
  - [Redpanda + kafdrop](#redpanda--kafdrop)
  - [Nats](#nats)
  - [Pubsub](#pubsub)

## Search

### Opensearch

```py
load('ext://helm_resource', 'helm_resource', 'helm_repo')
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
helm_resource(
    name='opensearch',
    chart='bitnami/opensearch',
    namespace='default',
    flags=[
        '--set=global.persistence.enabled=false',
        '--set=opensearchDashboards.enabled=false',
        '--set=replicas=1',
        '--set=opensearchSecurity.enabled=false',
    ],
    port_forwards=["9200:9200"],
    labels=['search']
)
```

## NoSQL

### [Cassandra](https://github.com/bitnami/charts/tree/main/bitnami/cassandra/)

TODO

## RMDBS

### Postgres

```py
load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://secret', 'secret_from_dict')
k8s_yaml(secret_from_dict("tilt-pg", inputs = {
    'postgres-password' : "password"
}))
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
helm_resource(
    name='pg',
    chart='bitnami/postgresql',
    namespace='default',
    flags=[
        '--set=image.tag=15-debian-11',
        '--set=global.postgresql.auth.existingSecret=tilt-pg',
        '--set=primary.persistence.enabled=false',
    ],
    labels=['database']
)
docker_build("pg-migrate-image", "migrations", dockerfile="tilt/extensions/pg/migrate.Dockerfile")
k8s_yaml("./tilt/extensions/pg/migrate.yaml")
k8s_resource("pg", port_forwards=["5432:5432"], labels=["database"])
k8s_resource("pg-migrate", resource_deps=["pg"], labels=["database"])
```

## Cache / Distributed Locks

### Redis

```py
load('ext://helm_resource', 'helm_resource', 'helm_repo')
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
helm_resource(
    name='redis',
    chart='bitnami/redis',
    namespace='default',
    flags=[
        '--set=cluster.enabled=false',
        '--set=usePassword=false',
        '--set=auth.enabled=false',
        '--set=master.persistence.enabled=false',
        '--set=slave.persistence.enabled=false',
        '--set=replica.replicaCount=1',
    ],
    port_forwards=["6379:6379"],
    labels=['cache']
)
```

## Message Queue

### Redpanda + kafdrop

[kafka](https://github.com/redpanda-data/helm-charts):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redpanda
  labels:
    app: redpanda
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redpanda
  template:
    metadata:
      labels:
        app: redpanda
    spec:
      enableServiceLinks: false
      containers:
      - name: redpanda
        image: docker.redpanda.com/vectorized/redpanda:latest
        ports:
        - containerPort: 9092
        - containerPort: 9644
        args:
        - redpanda
        - start
        - --smp 1
        - --reserve-memory 0M
        - --overprovisioned
        - --node-id 0
        - --check=false
        - --advertise-kafka-addr=redpanda:9092
        volumeMounts:
        - name: data
          mountPath: /var/lib/redpanda/data
      volumes:
      - name: data
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: redpanda
spec:
  selector:
    app: redpanda
  ports:
  - name: kafka
    port: 9092
    targetPort: 9092
  - name: admin
    port: 9644
    targetPort: 9644
```

[kafdrop](https://github.com/obsidiandynamics/kafdrop):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafdrop
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kafdrop
  template:
    metadata:
      labels:
        app: kafdrop
    spec:
      containers:
      - name: kafdrop
        image: obsidiandynamics/kafdrop
        ports:
        - containerPort: 9000
        env:
        - name: KAFKA_BROKERCONNECT
          value: "redpanda:9092"
---
apiVersion: v1
kind: Service
metadata:
  name: kafdrop
spec:
  ports:
  - port: 9000
    targetPort: 9000
```

### Nats

TODO

### Pubsub

TODO
