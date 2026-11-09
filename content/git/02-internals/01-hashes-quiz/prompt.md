# What's in a Commit Id?

Every commit has an id like `9a04b31c5e2d8f7a6b1c0e9d8f7a6b5c4d3e2f10`. It's not a random number. It's a **hash**: a 40-character fingerprint that Git calculates from the commit's contents.

A hash works like this:

- The same contents always give the same hash.
- Change **anything**, even one letter, and the hash changes completely.
- You can't work backwards from a hash to the contents.

What goes into a commit's hash?

- a snapshot of every file
- the parent commit's id
- the author, the date and time
- the message

You usually only need the first 7 characters (`9a04b31`). That's enough to be unique in almost any project.

(Git uses an algorithm called SHA-1 to calculate the hash. How it works is out of scope: just know ids are fingerprints.)

**You change one word in an old commit's message. What happens to its id?**
