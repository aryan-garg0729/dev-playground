# CheckList

The result panel under the editor after Run or Submit. It renders the grader's `Result`: a summary line and one row per check.

- **Summary**: `section` type. All passing: "3 of 3 checks passing" in `success`. Any failing: "2 of 3 checks failing" in `danger`. Runtime on the right in `caption`/`text-muted`.
- **Row**: an 8px status dot (`success` or `danger`, with `aria-label` "Passed"/"Failed"), then the call and its result in `code-sm`. Failed rows show the wrong value in `danger` followed by "(expected 6)" in `text-muted`. Never rely on the dot alone.
- Hidden checks show only their name and feedback message, never test source.

Sits in a `surface` pane with `radius-lg`, 14px × 16px padding, 6px between rows. The editor marks the offending line with a `danger-soft` background and a `danger` line number.
