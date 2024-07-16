# DNS Server and CLI Tool
This repository contains a DNS server and a command-line interface (CLI) tool for querying DNS records.
The server implemented using [DNS library](https://github.com/miekg/dns)supports both standard DNS and DNS over HTTPS (DoH), with features including caching, ad-blocking, and reverse DNS resolution.
The CLI tool allows users to perform DNS queries and reverse lookups with various options.

## Features
### DNS Server
* Standard DNS Support: Listens for DNS queries over UDP.
* DNS over HTTPS (DoH) Support: Listens for DoH queries over HTTPS.
* Caching: Reduces latency and load on upstream servers by caching DNS responses.
* Ad-Blocking: Blocks known ad domains using an open-source [ad-blocking list](https://github.com/StevenBlack/hosts).
* Reverse DNS Resolution: Supports PTR record queries for reverse DNS lookups.

### CLI Tool
* Query DNS Servers: Allows querying of DNS servers for various record types.
* DNS over HTTPS: Option to use DoH for queries.
* Reverse Lookups: Supports reverse DNS lookups for IP addresses.
* Verbose Output: Provides detailed output for DNS responses.
* Supported Record Types: Supports querying for multiple DNS record types (A, AAAA, MX, CNAME, PTR, etc.).

## Installation  
Clone the repository and navigate to the project directory:
```sh
git clone https://github.com/Harshitk-cp/dns-cli.git
cd dns-cli
```
Build the cli tool:
```sh
make build alias
```
Start server:
```sh
go run cmd/dns-server/main.go
```
and you're good to go!

> [!NOTE]
> In order to make use of DoH(DNS over HTTPS) feature you need to add path to the certificates in config.yaml.

```yaml
dnsServerAddr: ":53" # Standard DNS port
dohServerAddr: ":443" # Standard HTTPS port for DoH
dohCertFile: "path/to/cert.pem"
dohKeyFile: "path/to/key.pem"
```

You can generate your certificate in macOS using this command:
```sh
brew install mkcert
mkcert -install
mkcert -key-file ~/.cert/key.pem -cert-file ~/.cert/cert.pem "your IP"
```

### Installation using docker
Build the cli tool:
```sh
make build alias
```

Start server:
```sh
make docker-build && make run-server
```
> [!NOTE]
> Remember to add your path in the docker-compose file.
```yml
# docker-compose.yml

version: "3.8"

services:
  dns-server:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: dns-server
    ports:
      - "53:53/udp"
      - "443:443"
    volumes:
      - /path/.cert:/app/cert

```

## Using the CLI

### Basic Query
To query a DNS server for a domain:
```sh
dnscli query [server] [domain]
```

Exmaple:
```sh
dnscli query 8.8.8.8 google.com
# google.com.	84	IN	A	142.250.194.110
```
### Reverse Lookup
To perform a reverse DNS lookup:
```sh
dnscli query [server] [IP] --reverse
```

Exmaple:
```sh
dnscli query 8.8.8.8  8.8.8.8 --reverse
# 8.8.8.8.in-addr.arpa => dns.google.
```

### DNS over HTTPS
To query a DNS server using DoH:
```sh
dnscli query 192.168.1.14 netflix.com --doh
# netflix.com.	60	IN	A	54.155.178.5
# netflix.com.	60	IN	A	3.251.50.149
# netflix.com.	60	IN	A	54.74.73.31
```

### Verbose Output
To enable verbose output:
```sh
dnscli query 192.168.1.14 x.com --verbose
# ;; opcode: QUERY, status: NOERROR, id: 54567
# ;; flags: qr aa rd; QUERY: 1, ANSWER: 4, AUTHORITY: 8, ADDITIONAL: 0

# ;; QUESTION SECTION:
# ;x.com.	IN	 A

# ;; ANSWER SECTION:
# x.com.	1800	IN	A	104.244.42.129
# x.com.	1800	IN	A	104.244.42.1
# x.com.	1800	IN	A	104.244.42.65
# x.com.	1800	IN	A	104.244.42.193

# ;; AUTHORITY SECTION:
# x.com.	13999	IN	NS	d.u10.twtrdns.net.
# x.com.	13999	IN	NS	a.r10.twtrdns.net.
# x.com.	13999	IN	NS	a.u10.twtrdns.net.
# x.com.	13999	IN	NS	b.u10.twtrdns.net.
# x.com.	13999	IN	NS	b.r10.twtrdns.net.
# x.com.	13999	IN	NS	d.r10.twtrdns.net.
# x.com.	13999	IN	NS	c.u10.twtrdns.net.
# x.com.	13999	IN	NS	c.r10.twtrdns.net.
```

### Record Types
To specify a DNS record type:
```sh
dnscli query 192.168.1.14 www.linkedin.com --type AAAA # IPV6
# www.linkedin.com.	300	IN	CNAME	exp1.www.linkedin.com.
# exp1.www.linkedin.com.	300	IN	CNAME	www-linkedin-com.l-0005.l-msedge.net.
# www-linkedin-com.l-0005.l-msedge.net.	240	IN	CNAME	l-0005.l-msedge.net.
# l-0005.l-msedge.net.	240	IN	AAAA	2620:1ec:21::14
```

```sh
dnscli query 192.168.1.14 www.linkedin.com --type A # IPV4
# www.linkedin.com.	300	IN	CNAME	exp1.www.linkedin.com.
# exp1.www.linkedin.com.	300	IN	CNAME	www-linkedin-com.l-0005.l-msedge.net.
# www-linkedin-com.l-0005.l-msedge.net.	240	IN	CNAME	l-0005.l-msedge.net.
# l-0005.l-msedge.net.	240	IN	A	13.107.42.14
```

> [!TIP]
> You can even replace your current DNS with this one to keep things private. Running it locally ensures that only you have access to your data, putting your privacy entirely under your control.

Just put this url where required.
```sh
https://your IP/dns-query
```
