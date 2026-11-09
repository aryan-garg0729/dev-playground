# Tagging a Release

Most teams tag releases with **semantic versions**, `MAJOR.MINOR.PATCH`, and push the same image under several tags:

| Tag | Points at | Changes? |
| --- | --- | --- |
| `2.4.1` | Exactly this build | **Never** |
| `2.4` | The newest `2.4.x` | When a new patch comes out |
| `2` | The newest `2.x.y` | When a new minor or patch comes out |
| `latest` | The newest release overall | Every release |

Pushing several tags costs nothing: they all point at the same image, so the layers are uploaded once.

People then pick how adventurous they want to be. A careful production server pins `2.4.1` (or even the digest), and a developer happily runs `latest`.

In real teams this is automated: a **CI pipeline** (GitHub Actions, GitLab CI...) builds, tests, tags and pushes the image on every release, so nobody types these commands by hand.

You fixed a bug in version `2.4.0`, and the new release is `2.4.1`.

**What should you push?**
