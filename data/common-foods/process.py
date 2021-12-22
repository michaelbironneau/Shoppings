import json
import inflect
from string import digits, punctuation

p = inflect.engine()
foods = set()
common_words = set(["of", "a", "the", "and", "or"])
with open("./Food.json", "r", encoding="UTF-8") as f:
    for line in f:
        food = json.loads(line)
        food_name = food["name"].lower()
        food_name = food_name.translate(str.maketrans("", "", punctuation))
        food_name = food_name.translate(str.maketrans("", "", digits))
        for word in food_name.split(" "):
            if len(word) < 2 or word in common_words:
                continue
            foods.add(word)
            foods.add(p.plural_noun(word))

print(f"Read in {len(foods)} foods.")
f = open("./common-foods.json", "w")
f.write(json.dumps(list(foods)))
print("Done!")
