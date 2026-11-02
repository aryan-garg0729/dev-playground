```python
from datetime import date

prompt = ChatPromptTemplate.from_messages([
    ("system", "Today is {today}."),
    ("human", "{question}"),
]).partial(today=lambda: date.today().isoformat())
```

What does the `.partial(...)` call do?
