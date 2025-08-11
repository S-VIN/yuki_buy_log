# yuki_buy_log

## Docker Hub credentials

GitHub Actions build and push Docker images to Docker Hub. For this to work,
add the following secrets in the repository settings:

- `DOCKERHUB_USERNAME` – your Docker Hub username
- `DOCKERHUB_TOKEN` – a Docker Hub access token or password

Navigate to **Settings → Secrets and variables → Actions → New repository
secret** to create these entries.