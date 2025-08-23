# UniToken Development Instructions

## Architecture Overview

UniToken is a multi-component monorepo consisting of:

- **Go Service** (`service/`): Local API server providing token management and gateway functionality
- **Vue.js Frontend** (`frontend/`): Web UI for configuration and management
- **SDKs** (`sdk/`): Node.js and Python client libraries for service integration

## Technology Stack

- `pnpm` v10 for package management
- `go` v1.24 for backend service development
- `uv` for Python SDK development and service management
- `vue` v3 with setup API
- `shadcn-vue` for UI components. You may install components via `pnpx shadcn-vue add <component-name>`.
- `tailwindcss` v4 for styling

## i18n

- Always use English in the codebase.
- Use Vue I18n for frontend localization, with locale files in `locales/` (`en-US.json`, `zh-CN.json`)

## Code Style

- Prefer using `GET`, `POST`, instead of `PUT`, `DELETE` for API endpoints.
- Should extract common logic when possible.
- DO NOT expose any sensitive information in the codebase. This is a public repository.
- High code quality is expected, as this is a production-ready open-source project.
