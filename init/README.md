# init
scripts to run services

All of these services are planned to be run in a single podman pod.

## POD
```bash
podman pod create --name almagest
```

## REDIS
```bash
podman run -d --name almagest-redis -v $PWD/containers/redis/data:/var/redis/data:Z --pod almagest docker.io/library/redis:6.2-alpine
```

#### REDIS-CLI test pod
```bash
podman run --rm --name almagest-redis-cli --pod almagest -it docker.io/goodsmileduck/redis-cli
```


