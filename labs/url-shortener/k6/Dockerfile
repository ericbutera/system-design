FROM grafana/k6:latest
WORKDIR /app

COPY entrypoint.sh .
COPY *.js .

ENTRYPOINT [ "/app/entrypoint.sh" ]
