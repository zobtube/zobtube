# Docker guide

The easiest way to get started with Docker is by using the provided `docker-compose.yml` example.

## Docker Compose

1. Download the `docker-compose.yml` file [from this repo](https://raw.githubusercontent.com/zobtube/zobtube/master/docs/docker/docker-compose.yml)

    ```
    wget https://raw.githubusercontent.com/zobtube/zobtube/master/docs/docker/docker-compose.yml
    ```

2. Modify the `docker-compose.yml` file to your conveniance

    ```
    vim docker-compose.yml
    ```

3. Start Docker Compose

    ```
    docker compose up -d
    ```

4. ZobTube is now available on [http://localhost:8069](http://localhost:8069)

## Docker

If you prefer running Docker without the Compose part, you can use the following command:

```
docker run -it \
    -v ./zt-config:/config \
    -v ./zt-data:/var/lib/zobtube/data \
    -v ./zt-metadata:/var/lib/zobtube/metadata \
    -e ZT_DB_DRIVER=sqlite \
    -e ZT_DB_CONNSTRING=/config/db.sqlite \
    -e ZT_METADATA_TYPE=filesystem \
    -e ZT_METADATA_PATH=/var/lib/zobtube/metadata \
    -e ZT_MEDIA_PATH=/var/lib/zobtube/data \
    -e ZT_SERVER_BIND="0.0.0.0:8069" \
    -p 8069:8069 \
    ghcr.io/zobtube/zobtube
```

Then, ZobTube will be reachable on [http://127.0.0.1:8069](http://127.0.0.1:8069).
