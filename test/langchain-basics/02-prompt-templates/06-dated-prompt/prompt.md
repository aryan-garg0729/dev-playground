Write `make_dated_prompt()`. Build a prompt with the system message `Today is {today}.` and the human message `{question}`, then pre-fill `today` with `.partial()`.

Pass the `today` **function**, so the date is worked out each time the prompt is used, not once when it is built.
