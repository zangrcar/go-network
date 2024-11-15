from sentence_transformers import SentenceTransformer
import sys
import json


def process(str):
    model = SentenceTransformer('all-MiniLM-L6-v2')
    embedding = model.encode(str)
    return embedding

if __name__ == "__main__":
    string_input = sys.argv[1]
    result = process(string_input)
    print(json.dumps(result.tolist()))