# Fusili - Simple network port scanner

Fusili is a simple network port scanner. It is not a replacement for Nmap: it only performs TCP port scanning, and doesn't even try to be stealthy.

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