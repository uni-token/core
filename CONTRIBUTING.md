## Development

### Requirements

- Node.js 22
- PNPM 10
- go 1.24
- [uv](https://docs.astral.sh/uv/)

### Preparation

```sh
pnpm i
go -C service mod tidy
uv --directory sdk/python sync
```

### Start

```sh
pnpm -C ui dev
pnpm -C docs dev
go -C service run main.go debug
uv --directory sdk/python run example.py
```

### Build

```sh
pnpm -C ui build
pnpm -C docs build
go -C service build
```
