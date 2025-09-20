# Docker guide

The easiest way to get started with Docker is by using the provided `docker-compose.yml` example.

## Docker Compose

1. Download the `docker-compose.yml` file [from this repo](https://raw.githubusercontent.com/zobtube/zobtube/master/docs/docker/docker-compose.yml)
2. Modify the `docker-compose.yml` file to your conveniance
3. Open the terminal in the same directory
4. Start Docker Compose with `docker-compose up -d`
5. ZobTube is now available on [http://localhost:8069](http://localhost:8069)

## Docker

If you prefer running Docker without the Compose part, you can use the following command:

```
docker run -v ./zt-config:/config -v ./zt-data:/data -e ZT_DB_DRIVER=sqlite -e ZT_DB_CONNSTRING=/config/db.sqlite -e ZT_MEDIA_PATH=/data -e ZT_SERVER_BIND="0.0.0.0:8069" -p 8069:8069 ghcr.io/zobtube/zobtube
```

Then, ZobTube will be reachable on [http://127.0.0.1:8069](http://127.0.0.1:8069).
