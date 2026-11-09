# A Container Is a Process

Mo, our product manager, read that containers are "lightweight VMs". Run this and watch how fast it comes back:

```sh
time docker run --rm alpine:3.20 echo hi
# hi
# real  0m0.4s
```

A VM boots a whole operating system, which takes seconds. Docker skips the boot: it asks the kernel the host is *already running* for some walls:

- its own process list, so it thinks it's PID 1
- its own network interfaces and hostname
- its own filesystem, taken from an image

The walls are kernel features called **namespaces** (what it can see) and **cgroups** (how much CPU and memory it may use). How they work is out of scope for this course.

**Which description of a container is correct?**
