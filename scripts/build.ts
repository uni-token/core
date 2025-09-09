import process from 'node:process'
import { $ } from 'zx'

// Setup dependencies
await setupGo()
await setupUPX()

console.log('Building Go service...')
const ldflags = `-s -w -X 'logic.version=${new Date().toISOString().slice(0, 10).replace(/-/g, '')}' -X 'logic.appBaseUrl=http://uni-token.app'`
const releaseDir = 'ui/public/release'

await $`rm -rf ${releaseDir}`
await $`mkdir -p ${releaseDir}`

// Build for different platforms
await $`GOOS=linux GOARCH=amd64 go -C service build -ldflags=${ldflags} -o ../${releaseDir}/service-linux-amd64 ./main.go`
await $`GOOS=windows GOARCH=amd64 go -C service build -ldflags=${ldflags} -o ../${releaseDir}/service-windows-amd64.exe ./main.go`
await $`GOOS=darwin GOARCH=amd64 go -C service build -ldflags=${ldflags} -o ../${releaseDir}/service-darwin-amd64 ./main.go`

// Compress binaries with UPX
console.log('Compressing binaries...')
await $`upx --best --lzma ${releaseDir}/service-linux-amd64`
await $`upx --best --lzma ${releaseDir}/service-windows-amd64.exe`

// Build UI
console.log('Building UI...')
await $`pnpm -C ui build`

console.log('Build complete! Files are in ui/dist/')

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

async function setupUPX() {
  try {
    console.log('Checking UPX installation...')
    await $`upx --version`
    console.log('UPX is already installed.')
    return
  }
  catch {
    console.log('UPX is not installed, proceeding with installation...')
  }

  try {
    // Download and install UPX
    console.log('Downloading UPX...')
    await $`curl -LO https://github.com/upx/upx/releases/download/v5.0.2/upx-5.0.2-amd64_linux.tar.xz`
    await $`tar -xf upx-5.0.2-amd64_linux.tar.xz`
    await $`mv upx-5.0.2-amd64_linux/upx /usr/local/bin/`
    await $`rm -rf upx-5.0.2-amd64_linux upx-5.0.2-amd64_linux.tar.xz`

    // Verify installation
    await $`upx --version`
    console.log('UPX installation completed.')
  }
  catch (err: any) {
    console.log(`Failed to install UPX: ${err}`)
    throw err
  }
}
