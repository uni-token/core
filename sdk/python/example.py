from lib import request_unitoken_openai
from openai import OpenAI


def load_api_key() -> str | None: ...
def save_api_key(api_key: str | None): ...


def main():
    base_url, api_key = request_unitoken_openai(
        app_name="Example App",
        description="This is an example application.",
        saved_api_key=load_api_key(),
    )
    save_api_key(api_key)

    if not api_key:
        print("User rejected the request")
        return
    client = OpenAI(
        base_url=base_url,
        api_key=api_key,
    )
    demo_chat_streaming(client)


def demo_chat_streaming(client: OpenAI) -> None:
    stream = client.chat.completions.create(
        model="gpt-4o-mini",
        stream=True,
        messages=[
            {"role": "system", "content": "You are a concise assistant."},
            {"role": "user", "content": "Please write a one-sentence bedtime story."},
        ],
    )
    for chunk in stream:
        delta = chunk.choices[0].delta
        if getattr(delta, "content", None):
            print(delta.content, end="", flush=True)
    print()


if __name__ == "__main__":
    main()
