# Reading an Image Name

You've typed `alpine:3.20` and `nginx:1.27-alpine`. Those are short names. The full name of an image has up to three parts:

```
docker.io/library/nginx:1.27-alpine
└───┬───┘ └────┬─────┘ └────┬────┘
 registry   repository     tag
```

| Part | What it means | If you leave it out |
| --- | --- | --- |
| **registry** | Which online shop it lives in | Docker Hub (`docker.io`) |
| **repository** | The image's name, often `owner/name` | Required |
| **tag** | Which version | `latest` |

Official images on Docker Hub live under `library/`, which is why you can just say `nginx`. Images from people and companies look like `owner/name`, for example `grafana/grafana`. Images from other registries start with that registry's address, like `ghcr.io/...` (GitHub) or `quay.io/...`.

So `nginx:1.27-alpine` really means `docker.io/library/nginx:1.27-alpine`.

**How do you read `ghcr.io/acme/shop:2.1`?**
