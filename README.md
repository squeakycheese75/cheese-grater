# Cheese Grater: A Reverse Proxy and HTTP Redirector

Cheese Grater is a reverse proxy and HTTP redirector.  It's designed for one purpose, and that is to sit between an internet facing reverse proxy tool like [ngrok](https://ngrok.com/) or [localtunnel](https://theboroer.github.io/localtunnel-www/) and [LM Studio](https://lmstudio.ai/) and allow you to connect [AI Code Editor Cursor](https://www.cursor.com/) to use a local LLM, running on your local machine.

## Installation

Currently, you'll need [GO](https://go.dev/doc/install) installed.  

### Run with Go

1. Clone the repository:
    ```sh
    git clone https://github.com/squeakycheese75/cheese-grater.git
    ```

2. Navigate to the project directory:
    ```sh
    cd cheese-grater
    ```

3. Build the application:
    ```sh
    go build
    ```

## Install the ClI and run via the command line:

- Go to the [Releases](https://github.com/squeakycheese75/cheese-grater/releases) section on GitHub.
- Download the relevant binary for your platform, extract and run.
- If you want me to add a build for youR plaform, just raise an Issue   [Supported builds](https://golang.org/doc/install/source#environment)


### Usage - Run with go

1. Create a .env file:
    ```sh
    cp .env.sample .env
    ```

2. Update you .env file:
- APIKey: The API key for authentication.  You will need this for Cursor.
- RedirectURL: The URL of the backend server to which requests will be forwarded (default: localhost:1234).
- ProxyPort:  What port to run this app on.

### Usage - Platform CLI

Run the Cheese Grater application with the following command-line flags:

```sh
./cheese-grater -RedirectURL <backend-url> -Port <port> -APIKey <api-key>
```

### Command-Line Flags
- RedirectURL: The URL of the backend server to which requests will be forwarded (default: localhost:1234).
- Port: The port on which the Cheese Grater proxy will run (default: 8080).
- APIKey: The API key for authentication (default: generated at runtime).
- help: Display help information about command-line flags.


## How to connect Cursor to your local LLM.

Writing the article now...

## License
This project is licensed under the GNU License. See the LICENSE file for details.