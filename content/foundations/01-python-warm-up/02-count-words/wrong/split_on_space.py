def count_words(text):
    counts = {}
    for word in text.lower().split(" "):
        if word:
            counts[word] = counts.get(word, 0) + 1
    return counts
