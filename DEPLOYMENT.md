# Deployment Guide

This document describes the CI/CD pipeline and deployment options for the ChatApp backend.

## GitHub Actions CI/CD Pipeline

The project includes an automated CI/CD pipeline that builds and publishes Docker images to GitHub Container Registry (GHCR).

### Workflow Triggers

The workflow is triggered on:
- **Push to main/master branch**: Builds and pushes with `latest` tag
- **Push tags (v*)**: Builds and pushes with semantic version tags
- **Pull requests**: Builds only (no push to registry)

### Workflow Features

- **Multi-platform builds**: Supports both `linux/amd64` and `linux/arm64`
- **Build caching**: Uses GitHub Actions cache for faster builds
- **Security**: Uses GitHub's built-in `GITHUB_TOKEN` for authentication
- **Attestation**: Generates build provenance attestation for security
- **Smart tagging**: Automatic tag generation based on branch/tag/commit

### Generated Tags

| Trigger | Generated Tags |
|---------|----------------|
| Push to main | `latest`, `main`, `main-<commit-sha>` |
| Push tag v1.2.3 | `v1.2.3`, `1.2`, `1` |
| Push to feature branch | `<branch-name>`, `<branch-name>-<commit-sha>` |
| Pull request | No tags (build only) |

## Using the Docker Image

### Pull from GHCR

```bash
# Pull the latest version
docker pull ghcr.io/your-username/chatapp/chatapp-backend:latest

# Pull a specific version
docker pull ghcr.io/your-username/chatapp/chatapp-backend:v1.0.0
```

### Run the Container

```bash
# Basic run
docker run -p 8082:8082 ghcr.io/your-username/chatapp/chatapp-backend:latest

# With persistent data volume
docker run -d \
  --name chatapp-backend \
  -p 8082:8082 \
  -v chatapp-data:/root/data \
  ghcr.io/your-username/chatapp/chatapp-backend:latest

# With custom database path
docker run -d \
  --name chatapp-backend \
  -p 8082:8082 \
  -e DB_PATH=/app/data/chatapp.db \
  -v $(pwd)/data:/app/data \
  ghcr.io/your-username/chatapp/chatapp-backend:latest
```

## Setting Up GHCR Access

### For Repository Owners

1. **Enable Package Permissions**: Go to repository Settings → Actions → General → Workflow permissions
2. **Set to**: "Read and write permissions"
3. **Check**: "Allow GitHub Actions to create and approve pull requests"

### For Users Pulling Images

```bash
# Login to GHCR (if pulling private images)
echo $GITHUB_TOKEN | docker login ghcr.io -u your-username --password-stdin

# Or use a personal access token
echo $PAT | docker login ghcr.io -u your-username --password-stdin
```

## Local Development vs Production

### Local Development
- Use `docker-compose up --build` for local development
- Includes both backend and frontend services
- Uses local SQLite database
- Hot reload for development

### Production Deployment
- Use pre-built images from GHCR
- Configure external database (PostgreSQL recommended)
- Set up proper environment variables
- Use container orchestration (Docker Swarm, Kubernetes)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_PATH` | Database file path | `/root/chatapp.db` |
| `PORT` | Server port | `8082` |
| `JWT_SECRET` | JWT signing secret | Auto-generated |

## Monitoring and Logs

```bash
# View container logs
docker logs chatapp-backend

# Follow logs in real-time
docker logs -f chatapp-backend

# Check container status
docker ps
docker inspect chatapp-backend
```

## Troubleshooting

### Common Issues

1. **Permission denied when pushing to GHCR**
   - Check repository workflow permissions
   - Ensure `GITHUB_TOKEN` has package write permissions

2. **Build fails on multi-platform**
   - Check Dockerfile compatibility with ARM64
   - Ensure all dependencies support target platforms

3. **Container fails to start**
   - Check port conflicts (`docker ps`)
   - Verify volume mounts and permissions
   - Check environment variables

### Debug Commands

```bash
# Run container in interactive mode
docker run -it --rm ghcr.io/your-username/chatapp/chatapp-backend:latest sh

# Check image layers
docker history ghcr.io/your-username/chatapp/chatapp-backend:latest

# Inspect image details
docker inspect ghcr.io/your-username/chatapp/chatapp-backend:latest
```
