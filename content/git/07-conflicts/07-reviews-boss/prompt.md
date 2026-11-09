# The Reviews Page

Boss fight! Build a feature in messy commits while `main` changes under you, then land it cleanly with a real merge commit, the way a pull request would.

Order matters: squash **first**, rebase **second**. Three wip commits rebased over a conflict can mean resolving it three times; one squashed commit means once.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`add-reviews`** and make three wip commits. Copy and paste:

   ```sh
   printf '# Reviews\n\nWhat visitors say.\n' > reviews.md
   git add reviews.md && git commit -m "wip reviews"
   printf '\n## Leave a review\n\nEmail me!\n' >> reviews.md
   git commit -am "wip leave a review"
   sed 's/^Email: .*/Email: hello@example.com (replies within a day)/' contact.md > c.tmp && mv c.tmp contact.md
   git commit -am "wip reply time"
   ```

2. Meanwhile, the email address changes on `main`. Copy and paste:

   ```sh
   git switch main
   sed 's/^Email: .*/Email: hi@example.com/' contact.md > c.tmp && mv c.tmp contact.md
   git commit -am "H7: change my email address"
   ```

3. On `add-reviews`, squash the three wip commits into **`I7: add the reviews page`**.
4. Rebase it onto `main`. Resolve the conflict so the line is exactly:

   ```
   Email: hi@example.com (replies within a day)
   ```

5. Merge it into `main` as a **real merge commit** with the message **`J7: merge the reviews page`**. After the rebase, `main` *could* fast-forward; `git merge --no-ff` makes Git create a merge commit anyway.
6. Submit `gt.resolving.reviews-boss`.
