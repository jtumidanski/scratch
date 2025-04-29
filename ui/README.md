# Scratch Document UI

A document management service with a modern UI built with Next.js and shadcn/ui components.

## Features

- User authentication
- Document and folder management
- Real-time document editing
- Dark mode with Nord theme
- Responsive design

## Development

### Prerequisites

- Node.js 18 or later
- npm 8 or later

### Installation

```bash
# Install dependencies
npm install

# Start the development server
npm run dev
```

The application will be available at http://localhost:3000.

## Docker

### Building the Docker Image

```bash
# Build the Docker image
docker build -t scratch-document-ui .

# Run the Docker container
docker run -p 3000:3000 scratch-document-ui
```

The application will be available at http://localhost:3000.

### Environment Variables

The application uses the following environment variables:

- `NEXT_PUBLIC_API_BASE_URL`: The base URL of the API (default: http://localhost/v1)

You can set these variables when running the Docker container:

```bash
docker run -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL=http://api.example.com/v1 scratch-document-ui
```

## CI/CD

### GitHub Workflows

This repository includes two GitHub workflows:

1. **PR Validation**: Builds the Docker image and runs linting checks on pull requests to the master branch.
2. **Publish Docker Image**: Builds and publishes the Docker image to GitHub Packages when changes are merged into the master branch.

### Renovate

This repository uses Renovate to automatically update dependencies. The configuration includes:

- Automatic merging of minor and patch updates for stable dependencies
- Dependency dashboard for better visibility
- Scheduled updates on weekends to minimize disruption
- Docker image updates

## Deployment to k3s

### Prerequisites

- A running k3s cluster
- kubectl configured to access your cluster
- Helm (optional, for more complex deployments)

### Basic Deployment

1. Create a deployment YAML file:

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: scratch-document-ui
spec:
  replicas: 2
  selector:
    matchLabels:
      app: scratch-document-ui
  template:
    metadata:
      labels:
        app: scratch-document-ui
    spec:
      containers:
      - name: scratch-document-ui
        image: ghcr.io/your-username/scratch-document-ui:latest
        ports:
        - containerPort: 3000
        env:
        - name: NEXT_PUBLIC_API_BASE_URL
          value: http://api-service/v1
---
apiVersion: v1
kind: Service
metadata:
  name: scratch-document-ui-service
spec:
  selector:
    app: scratch-document-ui
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: scratch-document-ui-ingress
spec:
  rules:
  - host: scratch.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: scratch-document-ui-service
            port:
              number: 80
```

2. Apply the deployment:

```bash
kubectl apply -f deployment.yaml
```

### Using Helm (Optional)

For more complex deployments, you can create a Helm chart for the application.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
