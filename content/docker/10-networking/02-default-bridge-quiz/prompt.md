# Can You Hear Me?

Real apps are several containers working together: a website, an API, a database. They need to talk to each other. Let's try:

```sh
docker run -d --name cq-a cq/webapp:1.6
docker run --rm alpine:3.20 wget -qO- -T 3 http://cq-a:8080/
# wget: bad address 'cq-a:8080'
```

`bad address` means "I don't know any computer called `cq-a`". Both containers are on Docker's default network, and both are running.

**What's going on?**
