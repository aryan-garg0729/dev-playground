```python
prompt = ChatPromptTemplate.from_messages([
    ("system", "You translate {source} to {target}."),
    ("human", "{text}"),
])
```

What is `prompt.input_variables`?
