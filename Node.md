# Start a node
Base on the doc from official https://developer.algorand.org/docs/introduction-installing-node 
(testnet currently available)

## Installing on Docker

1. Pull latest code 
    ```bash
    git clone --branch rel/stable https://github.com/algorand/go-algorand.git
    ```
1. Build docker image
    ```bash
    cd docker/releases
    docker build -f Dockerfile-testnet . -t algorand/testnet-telem:latest
    ```
1. Create a Docker container
    * prepare configuration file
    ```bash
    docker run -it -v <root-path>:/root/node/tmp algorand/testnet
    node > cp -rf data tmp
    node > exit
    ```
    * copy configuration file `cp <root-path>/data/config.json.example <root-path>/data/config.json` from example and edit `vi data/config.json`, 
        ```bash
        "DNSBootstrapID": "testnet.algorand.network",
        ...

        "EndpointAddress": "127.0.0.1:8080",
        ```
    * start a container
    ```bash
    docker run -it -d -v <root-path>/data:/root/node/data --name algorand algorand/testnet
    ```
1. Sync Node with Network
    * run goal
        ```bash
        docker exec -it algorand /root/node/goal node start -d /root/node/data
        ```
    * print the rpc token
        ```bash
        docker exec -it algorand cat /root/node/data/algod.token
        ```
