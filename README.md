# analyzejs

This tools is an helper to get URLs from JS files.\
It can take either a list of urls as stdin or a folder containing JS files. \
It uses a common regex, [jsluice](https://github.com/BishopFox/jsluice) by Tomnomnom and BishopFox, and possibly a custom regex set by user. 

## Install

`go install github.com/0xRTH/analyzejs@latest`

## Usage : 

### Stdin:

`cat urls.txt | analyzejs`

### Folder

`analyzejs -???`

## Help : 


## Todos

- Finish readme
- Add flags

## Requirements

- Golang v1.18
