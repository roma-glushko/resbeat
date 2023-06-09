# 🔊resbeat

[![codecov](https://codecov.io/gh/roma-glushko/resbeat/branch/main/graph/badge.svg?token=BNJBL3XJ0O)](https://codecov.io/gh/roma-glushko/resbeat)

resbeat is a container agent that can expose container's resource usage via HTTP or websocket API:
- /ws/ - a websocket endpoint
- GET /usage/ - an HTTP polling endpoint

resbeat should be installed into the container's image and run along with the main container process.

resbeat could watch the following resources:
- general system resources via cgroup v1 (CPU and memory usage)

## Usage Report 

```
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
- general system resources via cgroup v2
- disk or volume utilization
- NVIDIA GPU utilization
