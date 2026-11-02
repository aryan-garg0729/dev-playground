from datetime import date

from cq import expect

from main import make_dated_prompt


def test_only_question_left():
    "only question is left to fill"
    expect(make_dated_prompt().input_variables, ["question"])


def test_formats_today():
    "system message has today's date"
    msgs = make_dated_prompt().format_messages(question="What day is it?")
    expect(msgs[0].content, f"Today is {date.today().isoformat()}.")
    expect(msgs[1].content, "What day is it?")
