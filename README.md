# Kiln & Co. – Go Application with Kubernetes, CI/CD and GitOps

Kiln & Co. is a small e-commerce website for a handmade ceramics shop, built with Go using only the standard library (`net/http` and `html/template`). It has no external Go dependencies or database.

This project goes beyond building the application. It demonstrates how to containerise a Go application, automate its build and deployment process with GitHub Actions, and deploy it to Amazon EKS using Helm and Argo CD.

The deployment workflow follows GitOps principles, with GitHub serving as the single source of truth for both application code and deployment configuration.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Technologies Used](#technologies-used)
- [Application Routes](#application-routes)
- [Project Structure](#project-structure)
- [Running Locally](#running-locally)
- [Running with Docker](#running-with-docker)
- [CI/CD Pipeline](#cicd-pipeline)
- [Deploying to Amazon EKS](#deploying-to-amazon-eks)
- [Configuring Argo CD](#configuring-argo-cd)
- [Accessing the Application](#accessing-the-application)
- [Cleaning Up](#cleaning-up)
- [Adding a Product](#adding-a-product)

## Features

- Landing page with featured products
- Shop catalogue and individual product pages
- Brand story page
- Enquiries form with a confirmation message
- Responsive layout using a single CSS file
- Docker containerisation using a multi-stage build
- Automated CI/CD pipeline with GitHub Actions
- Kubernetes deployment managed with Helm
- GitOps-based deployment using Argo CD

## Architecture

The project uses GitHub Actions for continuous integration and image publishing, while Argo CD handles continuous delivery through GitOps.

1. A developer pushes code changes to GitHub.
2. GitHub Actions builds the Go application and runs code-quality checks.
3. A Docker image is built and pushed to Docker Hub.
4. GitHub Actions updates the image tag in the Helm configuration and commits the change to GitHub.
5. Argo CD detects the updated configuration and synchronises the deployment with Amazon EKS.
6. Kubernetes pulls the new image and runs the application.

External traffic reaches the application through an AWS Load Balancer and NGINX Ingress Controller.

**Architecture Diagram**

<!-- Replace this placeholder with your architecture diagram -->

![Kiln & Co. Architecture Diagram](./docs/architecture.png)

## Technologies Used

| Technology | Purpose |
|---|---|
| Go | Application development |
| Docker | Application containerisation |
| Docker Hub | Container image registry |
| GitHub | Source code and deployment configuration |
| GitHub Actions | CI/CD automation |
| Amazon EKS | Managed Kubernetes cluster |
| Kubernetes | Container orchestration |
| Helm | Kubernetes package management |
| Argo CD | GitOps continuous delivery |
| NGINX Ingress | External application traffic routing |
| Cloudflare | DNS management |

## Application Routes

| Route | Description |
|---|---|
| `/` | Landing page |
| `/shop` | Product catalogue |
| `/shop/{id}` | Product detail |
| `/story` | About the brand |
| `/enquiries` | Contact and enquiry form |

## Project Structure

```text
kilnco/
├── .github/
│   └── workflows/
│       └── ci.yaml
├── helm/
│   └── go-web-app/
│       ├── templates/
│       │   ├── deployment.yml
│       │   ├── service.yaml
│       │   └── ingress.yml
│       └── values.yaml
├── static/
│   └── style.css
├── templates/
│   └── ...
├── main.go
├── go.mod
├── Dockerfile
└── README.md
```

## Running Locally

### Requirements

- Go 1.22 or newer
- Git

Clone the repository:

```bash
git clone https://github.com/Vivahjah/go-kubernetes-gitops-cicd.git
```

Navigate into the application directory:

```bash
cd go-kubernetes-gitops-cicd
```

Run the application:

```bash
go run .
```

The server starts on port `8080`.

Open the application in your browser:

```text
http://localhost:8080
```

## Running with Docker

### Requirements

- Docker installed and running

### 1. Build the Docker image

From the project root, run:

```bash
docker build -t go-web-app .
```

### 2. Run the container

```bash
docker run -d \
  --name kilnco \
  -p 8080:8080 \
  go-web-app
```

Open:

```text
http://localhost:8080
```

### 3. Check the running container

```bash
docker ps
```

To stop and remove the container:

```bash
docker stop kilnco
docker rm kilnco
```

## CI/CD Pipeline

The GitHub Actions workflow automates the build, validation and image publishing process.

The pipeline contains the following jobs:

1. **Build:** Compiles the Go application.
2. **Code Quality:** Runs `golangci-lint` to identify potential code-quality issues.
3. **Publish Docker Image:** Builds and pushes the Docker image to Docker Hub after the first two jobs succeed.
4. **Update Helm Values:** Updates the image tag in `helm/go-web-app/values.yaml` and commits the change back to GitHub.

### GitHub Repository Secrets

Before running the workflow, configure these repository secrets under:

`GitHub Repository → Settings → Secrets and variables → Actions`

| Secret | Description |
|---|---|
| `DOCKER_USERNAME` | Docker Hub username |
| `DOCKER_TOKEN` | Docker Hub access token with Read and Write permissions |
| `GH_PAT` | GitHub Personal Access Token with permission to update repository contents |

**Important:** Never commit access tokens, passwords or AWS credentials to the repository.

The workflow uses the GitHub Actions run ID as the Docker image tag, allowing each published image to have a unique version.

The Helm configuration is updated automatically after the image has been published.

## Deploying to Amazon EKS

This section explains how to recreate the Kubernetes environment.

### Prerequisites

Install and configure the following tools:

- AWS CLI
- `kubectl`
- `eksctl`
- Helm
- Docker

Configure your AWS credentials:

```bash
aws configure
```

Replace the example region below if you are using a different AWS region.

```bash
aws sts get-caller-identity
```

### 1. Create the EKS Cluster

The cluster used for this project was created with two `t3.medium` worker nodes.

```bash
eksctl create cluster \
  --name go-web-app-cluster \
  --region us-east-1 \
  --with-oidc \
  --vpc-nat-mode Disable \
  --node-type t3.medium \
  --nodes 2 \
  --node-volume-size 20
```

Wait for cluster creation to complete.

Configure `kubectl` to communicate with the cluster:

```bash
aws eks update-kubeconfig \
  --region us-east-1 \
  --name go-web-app-cluster
```

Check the nodes:

```bash
kubectl get nodes -o wide
```

All worker nodes should eventually show `Ready`.

### 2. Install NGINX Ingress Controller

Apply the AWS deployment manifest:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.15.1/deploy/static/provider/aws/deploy.yaml
```

Check the controller:

```bash
kubectl get pods -n ingress-nginx
```

Retrieve the external Load Balancer hostname:

```bash
kubectl get svc -n ingress-nginx
```

Look for the `EXTERNAL-IP` value of `ingress-nginx-controller`.

### 3. Install Argo CD

Create the namespace:

```bash
kubectl create namespace argocd
```

Install Argo CD:

```bash
kubectl apply --server-side -n argocd \
  -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Expose the Argo CD server through a Load Balancer:

```bash
kubectl patch svc argocd-server -n argocd \
  -p '{"spec": {"type": "LoadBalancer"}}'
```

Check the services:

```bash
kubectl get svc -n argocd
```

Retrieve the initial administrator password:

```bash
kubectl get secret argocd-initial-admin-secret \
  -n argocd \
  -o jsonpath="{.data.password}" | base64 -d
```

The default username is:

```text
admin
```

Keep the password private and change it after your initial login.

## Configuring Argo CD

Argo CD needs to know which Git repository and Helm chart to monitor.

Create a file named `argocd-application.yaml`:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: go-web-app
  namespace: argocd
spec:
  project: default

  source:
    repoURL: https://github.com/Vivahjah/go-kubernetes-gitops-cicd.git
    targetRevision: HEAD
    path: helm/go-web-app

  destination:
    server: https://kubernetes.default.svc
    namespace: default

  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

Apply the configuration:

```bash
kubectl apply -f argocd-application.yaml
```

This configures Argo CD to monitor the Helm chart in the repository and automatically synchronise deployment changes.

Check the application:

```bash
kubectl get applications -n argocd
```

Check the deployed resources:

```bash
kubectl get pods
kubectl get deployments
kubectl get services
kubectl get ingress
```

The application should eventually become `Synced` and `Healthy` in Argo CD.

## Accessing the Application

The application uses NGINX Ingress to route incoming traffic through the AWS Load Balancer.

To use your own domain:

1. Retrieve the NGINX Ingress Load Balancer hostname:

```bash
kubectl get svc -n ingress-nginx
```

2. Update the hostname in your DNS provider by creating a CNAME record.

Example:

| Field | Value |
|---|---|
| Type | CNAME |
| Name | `kilnco` |
| Target | `<AWS_LOAD_BALANCER_HOSTNAME>` |
| Proxy status | DNS Only |
| TTL | Auto |

3. Update the hostname in `helm/go-web-app/values.yaml` and the Ingress configuration to match your domain.

Example:

```yaml
ingress:
  enabled: true
  className: nginx
  hosts:
    - host: <YOUR_DOMAIN>
      paths:
        - path: /
          pathType: Prefix
```

4. Commit and push the changes to GitHub.

Once DNS has propagated and the Ingress is configured, visit:

```text
http://<YOUR_DOMAIN>
```

**Note:** This configuration exposes the application over HTTP. HTTPS/TLS is not configured in this project.

## Cleaning Up

AWS resources can continue to incur charges even when the application is not receiving traffic.

To remove the EKS cluster and its associated resources created by `eksctl`, run:

```bash
eksctl delete cluster \
  --name go-web-app-cluster \
  --region us-east-1
```

Check the AWS console afterwards for any remaining resources that may still incur charges, such as load balancers, EBS volumes and Elastic IP addresses.

## Adding a Product

To add a new product, update the `products` slice in `main.go`:

```go
{"glaze-cup", "Glaze Cup", "Mugs", "Short description here.", "#b5835a", 18},
```

The fields are:

- ID
- Name
- Category
- Description
- Colour
- Price

The product will automatically appear in the shop catalogue.

---

## Author

**[Your Name]**

- GitHub: [@Vivahjah](https://github.com/Vivahjah)
- Project Repository: [Kiln & Co. – Go Kubernetes GitOps CI/CD](https://github.com/Vivahjah/go-kubernetes-gitops-cicd)

## Acknowledgements

This project was inspired by the DevOps implementation tutorial by Abhishek Veeramalla.

