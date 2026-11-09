from main import make_dated_prompt


def test_runs_at_invoke_time():
    "works out the date each time the prompt is used"
    value = make_dated_prompt().partial_variables.get("today")
    assert callable(value), "pass the function itself to .partial(), not its result"
