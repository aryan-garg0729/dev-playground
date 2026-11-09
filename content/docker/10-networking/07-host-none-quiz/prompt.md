# The Other Network Types

You've met three network types. Here they are with one more:

| `--network` | What the container gets |
| --- | --- |
| *(nothing)* | The **default** bridge network: internet, but no name lookup |
| `cq-net` (your own) | A bridge network **with** name lookup, shared with chosen neighbours |
| `none` | No network at all |
| `host` | **No wall**: the container uses your computer's network directly. `-p` is ignored, and the app's port opens straight on your computer |

`host` is sometimes used for speed, but it removes one of the walls that makes containers safe, so avoid it unless you have a reason. (On a Mac it's even stranger: "host" means Docker's hidden Linux VM, not your Mac.)

**A container converts uploaded files and must never send data anywhere. Which network should it use?**
