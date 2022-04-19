socks5-connect
===

Connect tcp via socks5 and relay between stdio and the connection.

# Example

~/.ssh/config
```
Host somehost
    Hostname somehost.somedomain.example
    ProxyCommand /path/to/socks5-connect --socks5 socks-proxy:socks-port --dest %h:%p
```

## Requirements

 - go 1.18+

## License 

SEE LICENSE
