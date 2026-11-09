# Where Did the Orders Go?

Mo restarted the shop "the clean way" and every order disappeared. Try it yourself:

```sh
docker run -d --name eb-old -p 18070:8080 elixirbytes/brewd:1.6.0
curl -X POST 'localhost:18070/orders?potion=phoenix'
curl localhost:18070/orders
# phoenix

docker rm -f eb-old
docker run -d --name eb-old -p 18070:8080 elixirbytes/brewd:1.6.0
curl localhost:18070/orders
#                     <- empty
docker rm -f eb-old
```

brewd writes orders to `/data/orders.log` inside the container.

**Where was the phoenix order stored?**
