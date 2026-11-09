# Launch the News Page

Boss fight! Build a feature on its own branch while `main` keeps moving, then bring it home.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`launch/news`** and make two commits on it:
   - **`F3: add the news page`**, creating `news.md`:

     ```sh
     printf '# News\n\n- The website is live!\n' > news.md
     ```

   - **`G3: link news from projects`**: `echo '- [News](news.md)' >> projects.md`
2. Switch to `main`, add a question to the FAQ and commit it as **`H3: add a FAQ question`**:

   ```sh
   echo '- How often is the news updated?' >> faq.md
   ```

3. Merge `launch/news` into `main` with the message **`I3: merge the news page`**.
4. Delete `launch/news` with the safe delete.
5. Stay on `main` with nothing left uncommitted, then submit `gt.branching.news-boss`.
