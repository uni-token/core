import process from 'node:process'
import { $ } from 'zx'

// Build Go service
await setupGo()

console.log('Building Go service...')
const ldflags = `-s -w -X 'logic.version=${new Date().toISOString().slice(0, 10).replace(/-/g, '')}' -X 'logic.appBaseUrl=http://uni-token.app'`
const releaseDir = 'frontend/public/release'

await $`rm -rf ${releaseDir}`
await $`mkdir -p ${releaseDir}`

// Build for different platforms
await $`GOOS=linux GOARCH=amd64 go -C service build -ldflags=${ldflags} -o ../${releaseDir}/service-linux-amd64 ./main.go`
await $`GOOS=windows GOARCH=amd64 go -C service build -ldflags=${ldflags} -o ../${releaseDir}/service-windows-amd64.exe ./main.go`
await $`GOOS=darwin GOARCH=amd64 go -C service build -ldflags=${ldflags} -o ../${releaseDir}/service-darwin-amd64 ./main.go`

// Compress binaries with UPX (if available)
try {
  console.log('Compressing binaries...')
  await $`upx --best --lzma ${releaseDir}/service-linux-amd64`
  await $`upx --best --lzma ${releaseDir}/service-windows-amd64.exe`
}
catch (err: any) {
  console.log(`UPX not available, skipping compression. ${err}`)
}

// Build frontend
console.log('Building frontend...')
await $`pnpm -C frontend build`

console.log('Build complete! Files are in frontend/dist/')

async function setupGo() {
  try {
    console.log('Checking Go installation...')
    await $`go version`
    console.log('Go is already installed.')
    return
  }
  catch {
    console.log('Go is not installed, proceeding with installation...')
  }

  await $`curl -LO https://go.dev/dl/go1.25.0.linux-amd64.tar.gz`
  await $`rm -rf /usr/local/go`
  await $`tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz`

  // Also add to profile for future sessions
  await $`echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile`
  await $`source $HOME/.profile`

  // Set PATH for the current process
  process.env.PATH = `/usr/local/go/bin:${process.env.PATH}`

  // Verify installation
  await $`go version`
  console.log('Go installation completed.')
}
