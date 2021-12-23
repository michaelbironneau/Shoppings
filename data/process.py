import json
import inflect
from string import digits, punctuation

p = inflect.engine()
keywords = set()
items = set()
common_words = set(["of", "a", "the", "and", "or"])
files = ["./toiletries/Items.txt", "./household-items/Items.txt", "./clothes/Items.txt"]
def process_file(file_name):
    with open(file_name, "r", encoding="UTF-8") as f:
        for line in f:
            if "*" in line:
                continue  # Skip categories for now
            if len(line.strip()) == 0:
                continue 
            item_name = line.strip().lower()
            item_name = item_name.translate(str.maketrans("", "", punctuation))
            item_name = item_name.translate(str.maketrans("", "", digits))
            items.add(item_name)
            for word in item_name.split(" "):
                if len(word) < 2 or word in common_words:
                    continue
                keywords.add(word)
                keywords.add(p.plural_noun(word))

for file in files:
    process_file(file)
  
print(f"Read in {len(keywords)} keywords.")
f = open("./keywords.json", "w")
f.write(json.dumps(list(keywords)))
f.close()
print(f"Read in {len(items)} item phrases.")
f = open("./items.json", "w")
f.write(json.dumps(list(items)))
f.close()
print("Done!")
