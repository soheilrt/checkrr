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

## Using Environment Variables for Secrets

Instead of hardcoding sensitive information like API keys in the configuration file, you can use environment variables to securely manage secrets. This approach enhances security and makes it easier to manage configurations across different environments.

### Steps to Use Environment Variables

1. **Modify the Configuration File**:
   In the `config.yaml` file, remove the `api_key` field.
2. **Environment Variable Override**:
   The application will automatically override the `api_key` field in the configuration file with the corresponding environment variable. For example, the `api_key` for the `radarr` client will be replaced by the value of `API_KEY_RADARR`.
3. **Set Environment Variables**:
   Define the required environment variables in your system or container. For example:

    ```shell
    export API_KEY_RADARR="your_radarr_api_key"
    export API_KEY_SONARR="your_sonarr_api_key"
    ```

    If you are using Docker, you can pass these variables to the container using the `-e` flag:

    ```shell
    docker run -d \
      --name checkrr \
      --restart always \
      -v $(pwd)/config.yaml:/etc/checkrr/config.yaml:ro \
      -e API_KEY_RADARR="your_radarr_api_key" \
      -e API_KEY_SONARR="your_sonarr_api_key" \
      soheilrt/checkrr:latest
    ```

By following these steps, you can securely manage secrets without exposing them in your configuration files.

## Example of Logs

When running Checkrr, you should see logs similar to the following:

```shell
checkrr  | time="2024-11-18T00:59:06Z" level=info msg="Download is OK [ID: 75223438]: Reason: Average speed is 379.2KB/s"
checkrr  | time="2024-11-18T01:13:07Z" level=info msg="Download is OK [ID: 75223438]: Reason: Download status is not downloading"
```