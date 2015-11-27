# Fusili - Simple network port scanner

*Fusili* is a simple network port scanner. It is not a replacement for Nmap: it only performs TCP port scanning, and doesn't even try to be stealthy.

## Build

Requirements:

 * Go compiler (>= 1.3.0)
 * [gb](http://getgb.io/)

To build *Fusili*, execute `make` at the root of the sources. If no error occurred
during build, binary is available in the `bin/` directory.

## Configuration

The configuration file format is JSON:

`hosts`: hosts to scan and their ports expected to be open. Example:

```
"hosts": {
  "195.154.240.134": [ 22, 53 ],
  "62.210.248.118": [ 22 ],
  "62.210.248.119": [ ]
}
```

`output`: scan report output destinations. Supported types:

* `stdout`: output open ports found on the console.
* `s3`: upload the report as JSON file to an Amazon S3 bucket. Settings:
  * `access_key`
  * `secret_key`
  * `region`
  * `bucket`
  * `file_path`

Example:

```
"output": {
  "stdout": {
    "type": "stdout"
  },

  "s3": {
    "type": "s3",
    "access_key": "XXXXXXXXXXXXXXXXXXXX",
    "secret_key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "region": "eu-west",
    "bucket": "monitoring",
    "file_path": "portscan/report.json"
  }
}
```

## Usage

```
$ fusili -h
Usage: fusili [OPTIONS]

Options:
   -c  path to configuration file
   -h  display this help and exit
   -l  logging level (error, warning, notice, info, debug)
   -sC  scan concurrency
   -sP  scan ports range (START:END)
   -sT  scan port timeout (in seconds)
   -v  display version and exit
```

Example:

```
$ fusili -c conf.json -sC 100 -sP 1:100
2015/11/27 23:29:56.330393 WARNING: report: 195.154.240.134: found port 80/tcp open
2015/11/27 23:29:56.330562 WARNING: report: 62.210.248.119: found port 22/tcp open
2015/11/27 23:29:56.330588 INFO: core: scanned 3 hosts in 3.0 seconds
```
