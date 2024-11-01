# get-pcap-stats

A Go tool to collect various statistics from `.pcap` or `.pcapng` files using [tshark](https://tshark.dev/).

## Overview

This tool scans a directory for `.pcap` and `.pcapng` files and generates various statistical summaries for each file. The statistics are gathered using `tshark`, which allows for customizable and in-depth analysis of network data. Results are saved as text files with a specified suffix.

## Prerequisites

- [tshark](https://www.wireshark.org/docs/man-pages/tshark.html) must be installed and available in your system's PATH.

## Installation

Clone the repository and build the Go binary:

```bash
git clone https://github.com/yourusername/get-pcap-stats.git
cd get-pcap-stats
go build
```

## Usage

```bash
./get-pcap-stats [options]
```

### Options

| Option         | Description                                                                                             | Default                     |
|----------------|---------------------------------------------------------------------------------------------------------|-----------------------------|
| `-dir`         | Directory to search for `.pcap` or `.pcapng` files.                                                     | `.` (current directory)     |
| `-stats`       | Space-separated statistics types or path to a file containing stats types to pass to `tshark`.          | (all supported stats)       |
| `-suffix`      | Suffix for resulting statistics files.                                                                  | `.total-stats.txt`          |
| `-overwrite`   | If set, overwrites existing statistics files.                                                           | `false`                     |
| `-yes`         | If set, skips confirmation prompt before processing files.                                              | `false`                     |
| `-workers`     | Number of files to process in parallel.                                                                 | `NumCPU`                    |

### Examples

1. Process all `.pcap` files in the current directory and save results with `.summary.txt` suffix:
   ```bash
   ./get-pcap-stats -suffix ".summary.txt"
   ```

2. Process files in `/path/to/pcap` directory, overwriting any existing results:
   ```bash
   ./get-pcap-stats -dir /path/to/pcap -overwrite
   ```

3. Specify custom statistics types via a text file (one per line):
   ```bash
   ./get-pcap-stats -stats "path/to/stats_file.txt"
   ```

4. Find files in current directory and gather UDP statistics without confirmation prompt:
   ```bash
   ./get-pcap-stats -yes -stats "conv,udp endpoints,udp" -suffix ".stats-udp.txt"
   ```

## License

MIT License. See `LICENSE` for details.
