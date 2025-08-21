# UniToken

LLM Token Solution for Local AI Agents.

Open-source, local-first, and lightweight.

[**Docs**](https://docs.uni-token.app) | [Python SDK](https://docs.uni-token.app/sdk/python.html) | [Node.js SDK](https://docs.uni-token.app/sdk/node.html) | [Go SDK](https://docs.uni-token.app/sdk/go.html) | [Browser SDK](https://docs.uni-token.app/sdk/browser.html) | [Launch App](https://uni-token.ap) | [Development](./CONTRIBUTING.md)

## Why UniToken?

Local AI agents face a token management dilemma: developers must either run proxy servers and charge users for API costs, or force users to configure LLM tokens themselves, reducing usability.

UniToken provides a local solution. Through a lightweight SDK, AI agents can launch UniToken's interface for users to authorize existing tokens or purchase new ones from LLM providers - all on a single page. The agent then receives the Base URL and API Key.

UniToken consists of three open-source, local-first modules:

- **SDK**: Lightweight SDKs for Node.js and Python.
- **App**: A [static web interface](https://uni-token.app).
- **Service**: Local service handling storage and connections.

## Business Model

- UniToken is completely free - it simply guides users to LLM providers (SiliconFlow, OpenRouter, etc.), and shares the API Key with AI agents.
- Traditional AI software bundles token costs into subscriptions, creating problems: developers lose money as usage grows, and light users avoid high fees. UniToken solves this by letting users buy their own tokens and share them across local AI agents.
- We plan to establish partnerships with LLM providers, creating mutual benefits for providers, AI developers, and users.

## Project Architecture

<img src="https://docs.uni-token.app/arch.png" alt="UniToken Architecture" />

## Getting Started

For AI Agent developers, you may choose from the following lightweight SDKs to integrate UniToken into your AI agents:

- [Python SDK](https://docs.uni-token.app/sdk/python.html)
- [Node.js SDK](https://docs.uni-token.app/sdk/node.html)
- [Go SDK](https://docs.uni-token.app/sdk/go.html)
- [Browser SDK](https://docs.uni-token.app/sdk/browser.html)

For users, you can visit the [UniToken App](https://uni-token.app) to manage your LLM tokens.
