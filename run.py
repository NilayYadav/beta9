import os
import openai

client = openai.OpenAI(
  api_key="vhsYjoj7rBdJkEkN3qNX7IO4Bl8zK45Txw6VAA-iA4kNcgXpd4dip2VH9wShHQPT_7FWR7iAoCVUxpON7gdlEA==",
  base_url="https://llama3-1-8b-85c77e3-v1.app.beam.cloud/v1",
)

stream = client.chat.completions.create(
  model="Qwen/Qwen2.5-7B-Instruct-Turbo",
  messages=[
    {"role": "system", "content": "You are a travel agent. Be descriptive and helpful."},
    {"role": "user", "content": "Tell me about San Francisco"},
  ],
  stream=True,
)

for chunk in stream:
  print(chunk.choices[0].delta.content or "", end="", flush=True)