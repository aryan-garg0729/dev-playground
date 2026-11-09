# Air-Gapped Delivery

Priya's ticket:

> **STORE-114: storefront for the brewing server.** The brewing server is air-gapped (no internet, no registry). Ship it our web front end, `nginx:1.27-alpine`, under our own name. The ops team there checks every delivery against a manifest before loading it, so the manifest must name the exact bytes and how many layers to expect.

## Assignment

Work in `~/elixirbytes` and use only what this chapter covered.

1. Give `nginx:1.27-alpine` our own name: **`elixirbytes/storefront:1.0`**. No copies.
2. Create a **`delivery/`** folder and save **only** `elixirbytes/storefront:1.0` into **`delivery/storefront.tar`**. The archive must not carry the `nginx:1.27-alpine` tag: the brewing server should only ever see our name.
3. Write **`delivery/MANIFEST`** with exactly these lines:

   ```
   image=elixirbytes/storefront:1.0
   digest=sha256:<the image's 64-character digest>
   layers=<its number of filesystem layers>
   ```

4. Self-check: `tar -xOf delivery/storefront.tar manifest.json` should mention your tag and nothing else.

Submit from `~/elixirbytes` to attack. The checks compare your tag with the original image, read the archive's manifest, and look up the real digest and layer count to compare with `delivery/MANIFEST`.

Every check is a hit on the boss. Defeat it to unlock **Dockerfiles**.
