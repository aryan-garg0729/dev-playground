# ARG vs ENV

Dockerfiles have two kinds of variables, and they live at different times:

| | `ARG` | `ENV` |
| --- | --- | --- |
| Exists during | the **build** | the build **and** every container started from the image |
| Set with | `docker build --build-arg NAME=value` | the Dockerfile, or `docker run -e NAME=value` at run time |
| Default | `ARG NAME=default` | `ENV NAME=value` |

```dockerfile
FROM alpine:3.20
ARG SHOP=potions
ENV TIER=gold
RUN echo "building for $SHOP ($TIER)"
CMD ["sh", "-c", "echo shop=$SHOP tier=$TIER"]
```

**Which of these are true?** Pick every one that is.
