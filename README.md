# KVM

Key Vault Manager (KVM) is an app to help you find Azure 
Key Vaults, secrets, certificates and keys.

![KVM logo](logo.png)

## Install

### Go install

Run `go install` to download the binary to the go's binary folder:

```bash
go install github.com/nu12/kvm@latest
```

Note: go's binary folder (tipically `~/go/bin`) should be added to your PATH.

### From release

Download a tagged release binary for your OS (ubuntu, macos, windows) placing it in a folder in your PATH and make it executable (may require elevated permissions):

```bash
wget -O /usr/local/bin/kvm https://github.com/nu12/kvm/releases/download/<release>/kvm-<archtecture>.zip
chmod +x /usr/local/bin/kvm
```

### From source

Clone this repo and compile the source code:

```bash
git clone github.com/nu12/kvm
cd kvm
go build -o kvm main.go bundled.go
```

Move binary to a bin folder in your PATH (may require elevated permissions):
```bash
mv kvm /usr/local/bin/
```
