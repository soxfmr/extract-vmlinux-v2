# extract-vmlinux-v2
Golang implementation of [extract-vmlinux](https://github.com/torvalds/linux/blob/master/scripts/extract-vmlinux) script.

This project was forked from [https://github.com/Caesurus/extract-vmlinux-v2](https://github.com/Caesurus/extract-vmlinux-v2) and was completely reconstructed.

## Installing

```shell
go install github.com/soxfmr/extract-vmlinux-v2@latest
extract-vmlinux-v2 -file /boot/vmlinuz -output /tmp/vmlinux
```

## Usage

```shell
go get -u github.com/soxfmr/extract-vmlinux-v2@latest
```

decompress the kernel image and hold in memory:
```golang
file, err := os.Open("/boot/vmlinuz")
if err != nil {
    log.Fatalf("couldn't open the input file: %s", err)
}

kernel, err := vmlinux.Extract(file)
if err != nil {
    log.Fatalf("couldn't extract the kernel image: %s", err)
}
```

decompress the kernel image and save it to a file:
```golang
file, err := os.Open("/boot/vmlinuz")
if err != nil {
    log.Fatalf("couldn't open the input file: %s", err)
}

tmpFile, err := os.CreateTemp("", "vmlinux")

if err := vmlinux.ExtractTo(kernelFile, tmpFile); err != nil {
    log.Fatalf("couldn't extract the kernel image: %s", err)
}
```

---

## Compressions

### Supported
At the moment I have only found good native support for `gzip`/`bzip`/`lzma`/`lz4`/`xz`.

### Not Supported
I haven't found good support for `zstd`/`lzop`.

