# 🔊resbeat

[![codecov](https://codecov.io/gh/roma-glushko/resbeat/branch/main/graph/badge.svg?token=BNJBL3XJ0O)](https://codecov.io/gh/roma-glushko/resbeat)

resbeat is a container agent that can expose container's resource usage via HTTP or websocket API:
- `/ws/` - a websocket endpoint
- GET `/usage/` - an HTTP polling endpoint

resbeat should be installed into the container's image and run along with the main container process. 
Then, you should expose resbeat's port to let the rest of your system to scrape/consume container's/pod's utilization.
This is useful for building functionality around the usage reports like showing the user's env utilization somewhere in your UI.

resbeat could watch the following resources:

- general system resources via cgroup v1 or v2 (CPU and memory usage)
- NVIDIA GPU support

## Installation

```bash
curl -fSL https://github.com/roma-glushko/resbeat/releases/download/1.0.4-dev2/resbeat_Linux_x86_64.tar.gz -o "./resbeat_Linux_x86_64.tar.gz" \
    && tar -vxf resbeat_Linux_x86_64.tar.gz \
    && chmod +x ./resbeat
```

## Usage Report 

```json
{
    "collectedAt": "2023-06-11T20:01:49.851553Z",
    "cpu": {
        "usageInNanos": 150000000,
        "limitInCors": 3,
        "usagePercentage": 0
    },
    "memory": {
        "usagePercentage": 0.1220703125,
        "limitInBytes": 1073741824,
        "usageInBytes": 131072000
    }
}
```

## Plans

resbeat is intended to support more resource types like:
- disk or volume utilization
