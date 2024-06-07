# URL Shortener Microservice

This is a simple microservice written in Go for URL shortening. It provides APIs to create short URLs from long ones and redirect users to the original long URLs using the short ones.

## Features

- Shorten long URLs to short ones.
- Redirect users from short URLs to the original long URLs.

## Installation

### Prerequisites

- Go 1.16 or higher installed
- Redis

### Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/url-shortener.git
    cd url-shortener
    ```

2. Ensure redis is running 
    ```
    docker run --name redis -d -p 6379:6379 redis
    ```

3. Build the project:

    ```bash
    go build -o url-shortener main.go
    ```

4. Run the executable:

    ```bash
    go run main.go
    ```

## Usage

### API Endpoints

- `POST /shorten`: Shorten a long URL.
- `GET /{shortCode}`: Redirect to the original long URL associated with the provided short code.

### Example

```bash
curl -X POST -d '{"url": "https://example.com"}' http://localhost:8080/shorten

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
