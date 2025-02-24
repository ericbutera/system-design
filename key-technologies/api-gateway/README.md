# API Gateway

- Centralized API management
- Rate limit
- Authorization and authentication

## Implementations

- [AWS API Gateway](https://aws.amazon.com/api-gateway/)
- [GCP API Gateway](https://cloud.google.com/api-gateway)

## Experiments

KrakenD

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: krakend
  labels:
    app: krakend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: krakend
  template:
    metadata:
      labels:
        app: krakend
    spec:
      containers:
        - name: krakend
          image: devopsfaith/krakend
          args: ["run", "-c", "/etc/krakend/krakend.json"]
          ports:
            - containerPort: 8080
          volumeMounts:
            - name: config-volume
              mountPath: /etc/krakend
      volumes:
        - name: config-volume
          configMap:
            name: krakend-config

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: krakend-config
data:
  krakend.json: |
    {
      "version": 3,
      "name": "KrakenD API Gateway",
      "port": 8080,
      "endpoints": [
        {
          "endpoint": "/hotels",
          "method": "GET",
          "backend": [
            {
              "url_pattern": "/hotels",
              "host": ["http://hotel-service.default.svc.cluster.local:8000"]
            }
          ]
        },
        {
          "endpoint": "/reservations",
          "method": "GET",
          "backend": [
            {
              "url_pattern": "/reservations",
              "host": ["http://reservation-service.default.svc.cluster.local:8000"]
            }
          ]
        },
        {
          "endpoint": "/payments",
          "method": "GET",
          "backend": [
            {
              "url_pattern": "/payments",
              "host": ["http://payment-service.default.svc.cluster.local:8000"]
            }
          ]
        }
      ]
    }

---
apiVersion: v1
kind: Service
metadata:
  name: krakend
spec:
  selector:
    app: krakend
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080

```
