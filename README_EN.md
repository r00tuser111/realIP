# realIP

English | [简体中文](README.md)

An efficient command-line tool for filtering and identifying the actual IP addresses of target domains.

## Core Logic
- Filter domains that resolve to multiple IPs (supported by default)
- Filter domains that resolve to CDN CNAME records (parameter supported)
- Ultimate goal: only obtain domains with direct A record resolution to a single IP

## Features

- Support single domain or batch domain list resolution
- Support custom output formats (JSON, CSV)
- Support multi-threaded concurrent processing
- Provide comprehensive logging output
- Support custom output location
- Support CNAME record filtering, only showing domains that resolve directly to IP

## Installation

### Build from Source

```bash
git clone https://github.com/r00tuser111/realIP.git
cd realIP
go build -o realip
```

## Usage

### Command Line Arguments

- `-target`, `-u`: Specify a single target domain
- `-list`, `-l`: Specify a file path containing multiple domains (one per line)
- `-json`: Output results in JSON format
- `-csv`: Output results in CSV format
- `-threads`, `-t`: Set the number of concurrent threads (default: 10)
- `-o`: Specify output file path
- `-silent`: Disable all log output
- `-cname`: Filter CNAME records, only show domains that resolve directly to IP

### Examples

Resolve a single domain:
```bash
./realip -u example.com
```

Resolve a list of domains:
```bash
./realip -l domains.txt
```

Output in JSON format and save to file:
```bash
./realip -u example.com -json -o results.json
```

Output in CSV format and save to file:
```bash
./realip -l domains.txt -csv -o results.csv -t 20
```

Run in silent mode:
```bash
./realip -l domains.txt -o results.txt -silent
```

Filter CNAME records, only get domains that resolve directly to IP:
```bash
./realip -l domains.txt -cname -o results.txt
```

### Output Formats

#### Text Format (Default)
```
Domain          IP Address       CNAME    Status
----------------------------------------
example.com    93.184.216.34   No       Success
cdn.site.com   -               Yes      Failed
```

#### JSON Format
```json
[
  {
    "domain": "example.com",
    "ip": "93.184.216.34",
    "has_cname": false,
    "timestamp": "2024-03-18T15:04:05Z",
    "error": ""
  }
]
```

#### CSV Format
```csv
Domain,IP,Has_CNAME,Timestamp,Error
example.com,93.184.216.34,false,2024-03-18T15:04:05Z,
cdn.site.com,,true,2024-03-18T15:04:06Z,"Domain has CNAME record"
```

## License

MIT 