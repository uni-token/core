# UniToken

![Arch](https://docs.uni-token.app/arch.png)

<!--DEV-->

## Development

### Requirements

- Node.js 22
- PNPM 10
- go 1.24
- [uv](https://docs.astral.sh/uv/)

### Preparation

```sh
pnpm -C frontend i
go -C service mod tidy
uv --directory sdk/python sync
```

### Start

```sh
pnpm -C frontend dev
go -C service run main.go
uv --directory sdk/python run example.py
```

### Build

```sh
pnpm -C frontend build
go -C service build
```

<!--/DEV-->

## License

MIT LICENSE
