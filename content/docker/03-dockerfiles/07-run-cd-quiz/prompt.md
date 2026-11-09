# Spot the Bug

Mo found an old Dockerfile on the wiki and swears it "used to work":

```dockerfile
FROM golang:1.23-alpine
COPY app/ /src/
RUN cd /src
RUN go build -o /usr/local/bin/brewd .
CMD ["brewd"]
```

It fails on line 4.

**Why?**
