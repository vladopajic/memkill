# MemKill

Kills any process that exceeds specified memory.

## Why?

MemKill is a useful tool for managing processes with memory leak issues. It helps maintain system stability by monitoring memory usage and terminating processes that exceed a specified memory threshold. This proactive approach prevents the operating system from crashing or freezing due to excessive memory consumption.

## Usage

```sh
# memkill {max_memory_usage_in_MB}
memkill 100
```
This command runs `memkill`, which will kill any process whose memory usage exceeds 100 MB. The program continues running and monitoring processes until it is manually terminated.

## Running 
### Run using go

```sh
# install using go
go install github.com/vladopajic/memkill@latest

# run
$(go env GOPATH)/bin/memkill 100
```

### Run using downloaded binary
```sh
# install downloaded binary
wget https://github.com/vladopajic/memkill/releases/download/v0.0.1/memkill-linux-amd64
chmod +x memkill-linux-amd64
sudo mv memkill-linux-amd64 /usr/local/bin/memkill

# run
memkill 100
```

## Binaries

See latest binaries on [Releases](https://github.com/vladopajic/memkill/releases) page. 
Only linux binaries are currently available.
