# Service Scaling Experiment

This experiments shows performance characteristics of a service when scaled horizontally and vertically.

## Scenarios

### Single Instance

- set values.yaml `replicaCount` set to 1
- run k6 job
- observe number of requests

```mermaid
---
title: Single Instance
---
flowchart LR
  client <--> svc/api
    subgraph svc/api
      direction LR
      api1
    end
```

### Horizontal Scaling - Multiple Instances

Scaled Scenario:

- set values.yaml `replicaCount` values.yaml to 2, 5, 10, etc.
- ensure pods are ready (`kubectl get pods` should show 5 pods)
- run k6 job
- observe number of requests

```mermaid
---
title: Multiple Instances
---
flowchart LR
  client <--> svc/api
    subgraph svc/api
      direction LR
      api1
      api2
      api3
    end
```

#### Results

Resource limits:

```yaml
resources:
  requests:
    cpu: "50m"
    memory: "64Mi"
  limits:
    cpu: "100m"
    memory: "128Mi"
```

```mermaid
xychart-beta
  title "Horizontal Scaling Results"
  x-axis "Replicas" ["1", "5", "10", "15", "20"]
  y-axis "Total Requests" 15611 --> 402560
  bar [15611, 105233, 228752, 324160, 402560]
  line [15611, 105233, 228752, 324160, 402560]
```

### Vertical Scaling - Increase CPU Limits

TODO: add benchmarks

**NOTE**: The demo API doesn't perform cpu or memory intensive tasks, so horizontal scaling is more effective than vertical scaling.
