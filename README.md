# Checkrr: An Easy Way to Detect Dead/Stuck Downloads

Checkrr is a lightweight tool designed to help users detect stuck Radarr/Sonarr downloads and retry with an alternative
torrent. Built on a few simple ideas, Checkrr ensures reliability and value without introducing complications.

## **Looking for an example?** [HERE](./example)

## Checkrr Pillars

### Stateless Architecture

Checkrr doesn't keep track of any states internally. It relies solely on Sonarr/Radarr data to make decisions, allowing
seamless updates, restarts, or movements without affecting its behavior.

### Simple Stuck Detection Logic

Checkrr employs two fundamental concepts to identify stuck downloads:

1. **Average Download Speed Threshold**: This criterion helps identify both slow and dead downloads. The average
   download speed is calculated using a straightforward formula:
    ```plaintext
                                            (total file size - remaining download size) 
    average download speed per second = -------------------------------------------------
                                        number of seconds since the start of the download
    ```
   The advantages of this approach include:
    - Resilience to temporary network disconnections, as torrent is a distributed network.
    - The threshold and tolerance increase with the downloaded size and file size, allowing more time for reconnection
      or restarting the download for larger files.

2. **Download Timeout Threshold**: This concept provides a hard limit on download duration, ensuring downloads don’t
   take excessively long. Occasionally, files may download almost completely before getting stuck. In such cases, the
   average speed threshold might take a while to identify the issue. The timeout threshold is beneficial here, as it
   marks the download as stuck if it exceeds the specified duration.

## Installation

Running Checkrr is designed to be super easy. You only need two files to run it. You can find a working example of
Checkrr in the [example](./example) directory. You only need to modify the variables that are marked with `TODO` and
replace them with your values. After replacing the values, just run the following command:

```shell
docker compose up -d
```

1. **Configuration file**: This file holds all connection/criteria info of the application. You can find a complete
   example of the configuration file here ([example/config.yaml](example/config.yaml)).

2. **Running the docker image**: If you prefer to run Checkrr without docker compose, you can use the following command
   to spin up the application:

```shell
docker run -d \
  --name checkrr \
  --restart always \
  -v $(pwd)/config.yaml:/etc/checkrr/config.yaml:ro \ # TODO: Change the path to your config.yaml location
  soheilrt/checkrr:latest
```


## Example of Logs

When running Checkrr, you should see logs similar to the following:

```shell
checkrr  | time="2024-11-18T00:59:06Z" level=info msg="Download is OK [ID: 75223438]: Reason: Average speed is 379.2KB/s"
checkrr  | time="2024-11-18T01:13:07Z" level=info msg="Download is OK [ID: 75223438]: Reason: Download status is not downloading"
```