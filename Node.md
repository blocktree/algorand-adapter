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
    ```bash
    docker run -it -v <data-path>:/node/data -p 8080:8080 algorand/testnet
    ```
1. Sync Node with Network
    * copy configuration file `cp data/config.json.example data/config.json` from example and edit `vi data/config.json` (you might need install vim), 
        ```bash
        "DNSBootstrapID": "testnet.algorand.network",
        ...

        "EndpointAddress": "127.0.0.1:8080",
        ```
    * run goal
        ```bash
        ./goal node start -d data
        ```
    * print the rpc token
        ```bash
        cat data/algod.token
        ```
