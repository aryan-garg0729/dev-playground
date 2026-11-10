# Counting Under Symmetry

How many necklaces of n beads in m colours are there, if rotations of the same necklace count once? Counting all colourings and dividing by n is wrong, because some colourings look the same under several rotations.

**Burnside's lemma**: the number of distinct objects is the *average*, over all symmetries, of the number of colourings that the symmetry leaves unchanged.

For necklace rotations: rotating by k positions splits the beads into `gcd(k, n)` cycles, and a colouring stays the same exactly when each cycle is a single colour. So:

```text
distinct = (1 / n) · Σ_{k=0}^{n-1} m^gcd(k, n)
```

Modulo a prime, the `1 / n` is a multiplication by `pow(n, MOD - 2, MOD)`.

Recipe for any symmetry problem:

1. List the symmetries (n rotations; for a square grid, 4 rotations).
2. For each one, count the cycles it makes on the cells.
3. Add up `colours ^ cycles` and divide by the number of symmetries.

Reading: [Burnside's lemma on cp-algorithms](https://cp-algorithms.web.app/combinatorics/burnside.html) and the [video from your sheet](https://youtu.be/D0d9bYZ_qDY).

## Check yourself

A necklace has 4 beads in 2 colours, and rotations count as the same necklace. **How many distinct necklaces are there?**
