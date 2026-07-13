from elasticsearch import AsyncElasticsearch
import asyncio
import time


async def main():
    url = "http://localhost:9200"
    client = AsyncElasticsearch(url)

    if await client.ping():
        print("connected")
    else:
        print("❌ Connection failed. Is Docker running?")
        exit()

    index_name = "search_tutorial"

    if await client.indices.exists(index=index_name):
        print("Index already exists")

    print("Indexing 3 documents...")
    documents = [
        {"id": 1, "text": "Elastic search is incredibly fast"},
        {"id": 2, "text": "Search engines use an inverted index"},
        {"id": 3, "text": "Inverted index makes search fast"}
    ]

    for doc in documents:
        await client.index(index=index_name, id=doc.get("id"), document={"text": doc.get("text")})

    # Elasticsearch needs a split second to make new data searchable.
    # We force a refresh here for our script to work instantly.
    await client.indices.refresh(index=index_name)
    print("✅ Documents indexed successfully!\n")

    # Peeking under the hood: The Inverted Index (Analysis)
    # While Elasticsearch doesn't let you download the entire raw inverted index database
    # (it's heavily compressed binary data), you CAN ask it to show you exactly how it
    # breaks down text to build that index using the `analyze` API.
    print("--- HOW THE INVERTED INDEX SEES TEXT ---")
    sample_text = "Inverted index makes search fast!"
    print(f"Original Text: '{sample_text}'")

    analysis = await client.indices.analyze(text=sample_text)

    print("Extracted Tokens (These are what actually go into the Inverted Index):")
    for token in analysis["tokens"]:
        print(f"  - Token: '{token['token']}' (type: {token['type']})")

    print("\nNotice how it made everything lowercase and stripped the exclamation mark!\n")

    print("--- RUNNING A SEARCH ---")
    search_query = "fast search"
    print(f"Searching for: '{search_query}'")

    # We write a query telling Elasticsearch to look in the "text" field
    response = await client.search(
        index=index_name,
        query={
            "match": {
                "text": search_query
            }
        }
    )

    # Print the results
    hits = response["hits"]["hits"]
    print(f"Found {len(hits)} matching documents:")

    for hit in hits:
        doc_id = hit["_id"]
        score = hit["_score"]  # The relevance score (how well it matched)
        text = hit["_source"]["text"]
        print(f"  -> Document ID: {doc_id} | Score: {score} | Text: '{text}'")

if __name__ == "__main__":
    asyncio.run(main())
