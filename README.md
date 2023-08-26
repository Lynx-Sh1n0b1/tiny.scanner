# tiny.port Go scanner
```
# The lightweight simple scanner is written in Go for my own needs of process automation within Pentest.

Build package:
$ go build tiny.scanner.go
Run:                                                                                                                                                                                                                                                                                                                                          
$ ./tiny.scanner 
****************
* tiny.scanner *
****************

Usage: tiny.scanner -t <target> -p <ports> [-c <csv_file>]
  -c string
        CSV file to save results. Example -c output_results.csv
  -p string
        Port(s) to scan (e.g., 80,443 or 1-1024). Example -p 80 or -p 1-1024 or -p 22,80,3389,443
  -t string
        Target IPs (comma-separated), CIDR notation, or path to a text file with IPs per line. 
        Example -t 10.10.10.10 or -t 10.10.0.1,10.10.0.2 or -t 10.10.0.0/24 or -t ./ip_list.txt

                                                                                                                                                           
$ ./tiny.scanner -t your-target -p 1-1024
****************
* tiny.scanner *
****************
your-target 53 open
your-target 80 open
your-target 22 open
your-target 443 open

                                                                                                                                                   
$ ./tiny.scanner -t your-target -p 1-2000
****************
* tiny.scanner *
****************
your-target 53 open
your-target 80 open
your-target 22 open
your-target 443 open
your-target 2000 open
```
