# Host and None

Two special networks skip the bridge entirely:

| `--network` | The container's networking |
| --- | --- |
| `bridge` (default) or your own network | its own namespace and IP, NAT to the outside, `-p` to publish |
| `host` | **no isolation**: it shares the host's network stack |
| `none` | a loopback interface (`lo`) and nothing else |

```sh
docker run --rm --network none alpine:3.20 ip addr
# 1: lo: <LOOPBACK,UP,LOWER_UP> ...
#     inet 127.0.0.1/8 scope host lo
```

**Which of these are true?** Pick every one that is.
