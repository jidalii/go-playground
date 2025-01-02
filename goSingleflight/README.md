# Singleflight Example

- To run the example of cache miss with and without signle flight:

    ```shell
    cd cache_miss
    docker-compose up --build -d
    go test -bench=.
    ```
