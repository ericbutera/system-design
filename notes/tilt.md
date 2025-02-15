# Tilt notes

## Opensearch

```py
load('ext://helm_resource', 'helm_resource', 'helm_repo')
# helm_repo('opensearch-project', 'https://opensearch-project.github.io/helm-charts')
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
# k8s_resource("opensearch", port_forwards=["9200:9200"], labels=["search"])
```

## Postgres

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

## Redis

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
# k8s_resource("redis", port_forwards=["6379:6379"], labels=["cache"])
```
