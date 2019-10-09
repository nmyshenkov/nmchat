# nmchat - simple tcp chat server on golang

## Usage

For connect to server you can use Telnet

Example

```bash
telnet 127.0.0.1 3333

Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Write !help to show commands
 Your nikname: User1
```

## features

* Sending private messages to another clients in the chat
* Sending broadcast messages to all clients in the chat
* Changing nikname

## TODO list

* Redesign sending (send via nickname)
* Redesign chat window (attach input to bottom)
* Chat rooms
* Autocomplete
