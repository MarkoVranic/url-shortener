# url-shortener
Golang URL Shortener

Before proceeding, install the following prerequisites:

## Install

- [Docker](https://docs.docker.com/install/)
- [Docker compose](https://docs.docker.com/compose/install/)

Once everything is installed, execute the following commands from project root:

docker-compose up


## Usage

### Create short link
Input:
curl -X POST -H 'Content-Type: application/json' -d '{"url": "https://www.google.com"}' http://localhost/shorten

Output:
{"original_link":"https://www.google.com","short_link":"http://localhost/6mj9EwN"}

### Redirect to original link
curl  http://localhost/6mj9EwN -v
