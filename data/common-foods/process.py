import json
import inflect
from string import digits, punctuation

p = inflect.engine()
keywords = set()
foods = set()
common_words = set(["of", "a", "the", "and", "or"])
with open("./Food.json", "r", encoding="UTF-8") as f:
    for line in f:
        food = json.loads(line)
        food_name = food["name"].lower()
        food_name = food_name.translate(str.maketrans("", "", punctuation))
        food_name = food_name.translate(str.maketrans("", "", digits))
        foods.add(food_name)
        for word in food_name.split(" "):
            if len(word) < 2 or word in common_words:
                continue
            keywords.add(word)
            keywords.add(p.plural_noun(word))

print(f"Read in {len(keywords)} keywords.")
f = open("../../app/shoppings/src/app/shared/data/food-keywords.json", "w")
f.write(json.dumps(list(keywords)))
f.close()
print(f"Read in {len(foods)} food phrases.")
f = open("../../app/shoppings/src/app/shared/data/common-foods.json", "w")
f.write(json.dumps(list(foods)))
f.close()
print("Done!")
