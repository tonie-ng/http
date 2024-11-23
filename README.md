# blip
A simple HTTP server written in golang

## Installation
To install and run the server, you'll need to have Go installed on your machine.
- Clone the repository:
```bash
git clone https://github.com/yourusername/blip.git && cd blip
```
- Build the project:
```bash
go build -o blip
```
- Run the server:
```bash
./blip
```
By default, the server runs on `http://localhost:6703`.

## Usage
Once the server is running, you can interact with it using any HTTP client (browser, curl, Postman, etc.).
> For now the server has only one endpoint which is used to retrieve a html
> file

To get the HTML file, simply make a GET request to the root endpoint:
```bash
curl http://localhost:6703
```
If you want to only retrieve the headers (without the body), you can use the -I flag with curl:
```bash
curl -I http://localhost:6703
```
