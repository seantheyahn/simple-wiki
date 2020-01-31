# Simple Wiki
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/seanyahn/simplewiki) ![Docker Pulls](https://img.shields.io/docker/pulls/seanyahn/simplewiki)
## Usage
- Make sure you have [PostgreSQL] and [Go] installed on your system.
- Create a user and database owned by that user in PostgreSQL.
- Clone the repository, copy the config.json to your desired path and configure it to your need, don't forget to configure the ConnectionURI for the database.
- Go to to the project directory and run:
- ```sh
    go build -o simplewiki .
    ./simplewiki -config=/path/to/config.json
    ```
- Open the url (e.g. localhost:8080) in your browser and login with the root user which is configured in the config.json.

*(Note: The root user is created on the first run or whenever the user with id \<root\> does not exist in the DB, therefore changing it's username or password in the config won't update it after it is created.)*
### Config
Example config.json:
```json
{
    "Server": {
        "ListenAddress": ":8080",
        "CookieSecret": "my-super-secret-cookie-secret",
        "CSRFSecret": "my-csrf-token-secret",
        "DevelopmentMode": true,
        "IPRateLimiter": {
            "ReqPerSecond": 10,
            "ReqBurst": 10
        }
    },
    "DB": {
        "ConnectionURI": "postgres://test:test@localhost:5432/test?sslmode=disable",
        "MaxOpenConns": 25,
        "MaxIdleConns": 1
    },
    "RootUser": {
        "Username": "root",
        "Password": "my-complex-password"
    },
    "Logging": {
        "Level": "debug"
    }
}
```
You can override each level of the config file using environment variables, to do so you must prefix the variable name with ```CONF__``` and format the value as JSON. Each double-underscore ```__``` goes one level deeper in the JSON.
example:
```sh
export CONF__Server__ListenAddress='":8080"' #note the "" for JSON string
export CONF__Server__DevelopmentMode='true' #this is a JSON bool, no need for ""
export CONF__DB__ConnectionURI='"postgres://wiki:my-db-password@localhost:5432/wiki?sslmode=disable"'
```
## Docker Support
The image is available on [Docker Hub] as [seanyahn/simplewiki](https://hub.docker.com/r/seanyahn/simplewiki).
To quickly get it up and running using docker you can use the [docker-compose.yml](/docker-compose.yml).

License
----
MIT

[//]:#
   [Docker Hub]: <https://dockerhub.com>
   [PostgreSQL]: <https://postgresql.org>
   [Go]: <https://golang.org>
