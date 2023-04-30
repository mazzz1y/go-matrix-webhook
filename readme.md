# Go-Matrix-Webhook

Go-Matrix-Webhook is a lightweight and easy-to-use Matrix bot that allows you to send messages to a Matrix room using a simple POST request. It's perfect for integrating with other services, automating notifications, or even building your own custom applications.

## Installation

### Docker

Run the Docker container:

```bash
docker run -d --name go-matrix-webhook \
  -p 8080:8080 \
  -e LISTEN_ADDR="0.0.0.0" \
  -e LISTEN_PORT="8080" \
  -e LISTEN_PATH="/" \
  -e SECRET_HEADER="your_secret_header" \
  -e MATRIX_ACCESS_TOKEN="your_matrix_access_token" \
  -e MATRIX_ID="your_matrix_id" \
  -e MATRIX_URL="your_matrix_url" \
  ghcr.io/mazzz1y/go-matrix-webhook:latest
```

Replace `your_secret_header`, `your_matrix_access_token`, `your_matrix_id`, and `your_matrix_url` with the appropriate values.

### Download from GitHub Releases

1. Download and extract the latest release binary for your platform from the [GitHub Releases](https://github.com/mazzz1y/go-matrix-webhook/releases) page.

2. Run the binary:

```bash
./go-matrix-webhook
```

## Configuration

You can configure Go-Matrix-Webhook using the following environment variables:

- `LISTEN_ADDR`: The address on which the webhook server listens (default: "0.0.0.0").
- `LISTEN_PORT`: The port on which the webhook server listens (default: 8080).
- `LISTEN_PATH`: The path on which the webhook server listens (default: "/").
- `SECRET_HEADER`: Pre-shared `X-Secret` security header.
- `MATRIX_ACCESS_TOKEN`: The access token for your Matrix account.
- `MATRIX_ID`: Your Matrix user ID.
- `MATRIX_URL`: The URL of the Matrix homeserver.

Alternatively, you can set these options using command-line flags. Run `./go-matrix-webhook --help` for more information.

## Usage

### Create access token
```bash
curl -X POST -H "Content-Type: application/json" -d '{"type": "m.login.password", "identifier": {"type": "m.id.user", "user": "<bot username>"}, "password": "<bot password>", "initial_device_display_name": "Webhook Client"}' "https://<your-server>/_matrix/client/r0/login"
```
Get `access_token` from curl output. Don't forget to add bot to your room!
### Send Messages

To send a message, simply make a POST request with the following structure:

```bash
curl -X POST -H 'X-Secret: <your_secret_header>' -H "Content-Type: application/json" -d '{"message": "<your_message>", "room_id": "<your_room_id>"}' "http://<your_webhook_server>:<your_webhook_port>"
```

Replace `<your_secret_header>`, `<your_message>`, `<your_room_id>`, `<your_webhook_server>`, and `<your_webhook_port>` with the appropriate values.
