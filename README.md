# Tool for querying GeoIP info

This tool allows to add basic GeoIP info (country and city) to file with IP addresses, such as HTTP server access log.

To use you need to provide the path to GeoIP2 City database, you can get free one from https://dev.maxmind.com/geoip/geolite2-free-geolocation-data/#accessing-geolite2-free-geolocation-data.

The tool reads lines from standard input, it assumes that each line contains space-separated fields so you need to specify the field number for the IP address. It outputs the original line with two fields prepended - country and city.

Example usage:
```
tail -f access.log | ./geoip GeoIP2-City.mmdb 1
cat access.log | ./geoip GeoIP2-City.mmdb 1 | awk '{print $1 "/" $2}' | sort | uniq -c
cat access.log | ./geoip GeoIP2-City.mmdb 1 | awk '{print $1}' | sort | uniq -c | sort -nr
```