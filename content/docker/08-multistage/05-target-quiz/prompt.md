# Stop at a Stage

The build stage is thrown away, which is exactly what you want for shipping. But sometimes the build itself breaks, and you want to look around **inside the build stage**: is the code there? Does `go` work?

`docker build` can stop at any named stage and turn *that* into the image. It's also how teams run tests in a stage that has all the tools, without shipping those tools.

Run `docker build --help` and look for an option that sets "the target build stage".

**Which command builds only the stage named `build`, as `cq/webapp:debug`?**
