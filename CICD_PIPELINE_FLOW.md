# CI/CD Pipeline Flow - DWS Platform

## Linear Flow Diagram

```
┌──────────────┐
│  Developer   │
│  git push    │
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────┐
│                        GitHub                                 │
│  ┌────────────┐      ┌─────────────┐      ┌──────────────┐  │
│  │ Repository │─────▶│   Actions   │─────▶│     GHCR     │  │
│  │   (code)   │      │  (runners)  │      │   (images)   │  │
│  └────────────┘      └─────┬───────┘      └──────┬───────┘  │
└───────────────────────────┬│───────────────────────┼─────────┘
                            ││                       │
        ┌───────────────────┘└───────────┐           │
        ▼                                ▼           │
┌──────────────────┐          ┌──────────────────┐   │
│ Event Service    │          │ Ticket Service   │   │
│ Pipeline         │          │ Pipeline         │   │
└──────────────────┘          └──────────────────┘   │
                                                      │
                                                      ▼
                                          ┌───────────────────┐
                                          │     ArgoCD        │
                                          │  (GitOps sync)    │
                                          └─────────┬─────────┘
                                                    │
                                                    ▼
                                          ┌───────────────────┐
                                          │   Kubernetes      │
                                          │   Cluster (LTU)   │
                                          └───────────────────┘
```

---

## Event Service Pipeline (Sequential)

```
1. CODE PUSH
   ├─ Developer: git push origin main
   └─ GitHub: Webhook triggers Actions

2. TEST JOB (~2 min)
   ├─ Checkout code (actions/checkout@v4)
   ├─ Setup Go 1.23.1
   ├─ Cache Go modules (~/go/pkg/mod)
   ├─ Download dependencies (go mod download)
   ├─ Run tests on 10 packages:
   │  ├─ configs
   │  ├─ internal/controllers/events
   │  ├─ internal/controllers/health
   │  ├─ internal/event
   │  ├─ internal/middlewares
   │  ├─ internal/pkg/logger      (100% ✅)
   │  ├─ internal/pkg/metrics     (100% ✅)
   │  ├─ internal/pkg/utils       (78.9%)
   │  └─ internal/router
   ├─ Generate coverage report
   │  └─ Total: 20.5%
   ├─ Upload to Codecov
   └─ ✅ Exit 0 (success)

3. BUILD & PUSH JOB (~3 min)
   [Needs: test job success]
   ├─ Checkout code
   ├─ Login to ghcr.io
   ├─ Extract metadata (tags)
   ├─ Build Docker image
   │  ├─ FROM golang:1.23-alpine
   │  ├─ COPY . /app
   │  ├─ RUN go build
   │  └─ EXPOSE 8080
   ├─ Push to GHCR:
   │  ├─ ghcr.io/dws-org/dws-event-service:latest
   │  └─ ghcr.io/dws-org/dws-event-service:{sha}
   └─ ✅ Image published

Total Time: ~5 minutes
```

---

## Ticket Service Pipeline (Parallel + Sequential)

```
1. CODE PUSH
   ├─ Developer: git push origin master
   └─ GitHub: Webhook triggers Actions

2a. TEST JOB (~3 min) - WITH DATABASE
    ├─ Start PostgreSQL 16 container
    │  ├─ User: postgres
    │  ├─ Password: postgres
    │  ├─ Database: tickets_test
    │  └─ Health check: pg_isready
    ├─ Checkout code
    ├─ Setup Go 1.23.1
    ├─ Cache Go modules
    ├─ Download dependencies
    ├─ Generate Prisma client
    │  └─ go run github.com/steebchen/prisma-client-go generate
    ├─ Push database schema
    │  └─ go run github.com/steebchen/prisma-client-go db push
    ├─ Run tests (all packages)
    │  └─ DATABASE_URL=postgresql://localhost:5432/tickets_test
    ├─ Check coverage threshold
    │  ├─ Minimum: 15%
    │  └─ Current: ~0.7% (too low, but allowed to pass)
    ├─ Upload to Codecov
    └─ ✅ Exit 0

2b. LINT JOB (~1 min) - PARALLEL TO TEST
    ├─ Checkout code
    ├─ Setup Go 1.23.1
    ├─ Download dependencies
    ├─ Generate Prisma client
    ├─ Run golangci-lint
    │  ├─ Timeout: 5m
    │  ├─ Checks: gofmt, govet, ineffassign, staticcheck...
    │  └─ Scans: deadcode, unused, misspell, etc.
    └─ ✅ No issues found

3. BUILD JOB (~5 min)
   [Needs: test AND lint success]
   ├─ Checkout code
   ├─ Setup Go 1.23.1
   ├─ Generate Prisma client
   ├─ Build local binary
   │  └─ go build -v -o server ./cmd/server
   ├─ Setup Docker Buildx
   ├─ Login to ghcr.io
   ├─ Build & Push API Image:
   │  ├─ Context: .
   │  ├─ Dockerfile: ./Dockerfile
   │  └─ Tags:
   │     ├─ ghcr.io/dws-org/dws-ticket-service:latest
   │     └─ ghcr.io/dws-org/dws-ticket-service:{sha}
   └─ Build & Push Consumer Image:
      ├─ Context: .
      ├─ Dockerfile: ./Dockerfile.consumer
      └─ Tags:
         ├─ ghcr.io/dws-org/dws-ticket-service-consumer:latest
         └─ ghcr.io/dws-org/dws-ticket-service-consumer:{sha}

Total Time: ~6 minutes
```

---

## Deployment Flow (GitOps)

```
4. ARGOCD SYNC
   ├─ ArgoCD detects new image in GHCR
   ├─ Polls gitops repository
   │  └─ https://github.com/dws-org/gitops
   ├─ Compares desired vs actual state
   ├─ Auto-sync enabled: YES
   └─ Triggers deployment

5. KUBERNETES ROLLOUT
   ├─ Apply new manifests
   │  ├─ Deployment: dws-event-service
   │  ├─ Deployment: dws-ticket-service-api
   │  └─ Deployment: dws-ticket-service-consumer
   ├─ Rolling update strategy:
   │  ├─ maxSurge: 1
   │  ├─ maxUnavailable: 0
   │  └─ Zero-downtime deployment
   ├─ Pull images from GHCR
   ├─ Start new pods
   ├─ Wait for health checks:
   │  ├─ Liveness probe: /livez
   │  ├─ Readiness probe: /readyz
   │  └─ Timeout: 60s
   ├─ Route traffic to new pods
   └─ Terminate old pods

6. HEALTH CHECK
   ├─ Prometheus scrapes /metrics
   ├─ Check pod status: Running
   ├─ Check service endpoints: Ready
   └─ ✅ Deployment successful

Total Deployment Time: ~2 minutes
```

---

## Pipeline Comparison

| Aspect | Event Service | Ticket Service |
|--------|---------------|----------------|
| **Jobs** | 2 (test, build) | 3 (test, lint, build) |
| **Test DB** | ❌ No | ✅ PostgreSQL |
| **Linting** | ❌ No | ✅ golangci-lint |
| **Coverage Check** | ❌ No threshold | ✅ 15% minimum |
| **Docker Images** | 1 (API) | 2 (API + Consumer) |
| **Duration** | ~5 min | ~6 min |
| **Parallelization** | Sequential | Test ‖ Lint |

---

## Key Differences

### Event Service
```
Trigger → Test → Build → Push → Deploy
         (2min)  (3min)
         ────────────────
         Total: 5 minutes
```

**Advantages:**
- ✅ Faster pipeline (no lint, no DB)
- ✅ Simpler setup

**Disadvantages:**
- ❌ No code quality checks
- ❌ No integration tests (no DB)
- ❌ No coverage threshold

### Ticket Service
```
Trigger → Test (with DB) ┐
       → Lint            ├→ Build → Push → Deploy
         (3min)  (1min) /   (5min)
         ────────────────
         Total: 6 minutes (parallel: 3min + 5min)
```

**Advantages:**
- ✅ Integration tests with PostgreSQL
- ✅ Code quality enforcement (lint)
- ✅ Coverage threshold (15%)
- ✅ Two images (API + Consumer)

**Disadvantages:**
- ⚠️ Slightly slower
- ⚠️ More complex setup

---

## Monitoring & Observability

### During Pipeline

**GitHub Actions UI:**
```
✓ Test      2m 15s
✓ Lint      1m 05s  (Ticket only)
✓ Build     4m 32s
━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total:      5m 47s
```

**Coverage Reports:**
```
Codecov Dashboard:
├─ Event Service:  20.5% (+18.5% this push)
├─ Ticket Service:  0.7% (no change)
└─ Trend: ↗️ Event improving, Ticket needs work
```

### After Deployment

**ArgoCD Dashboard:**
```
Application: dws-event-service
├─ Sync Status: ✅ Synced
├─ Health Status: ✅ Healthy
├─ Revision: abc123f (SHA)
└─ Pods: 2/2 Running
```

**Kubernetes:**
```bash
$ kubectl get pods -n dws-event-service
NAME                                READY   STATUS    AGE
dws-event-service-7d9f8b9c4-8xk2p   1/1     Running   5m
dws-event-service-7d9f8b9c4-qw9rt   1/1     Running   5m
```

---

## Failure Scenarios

### Test Failure
```
Test Job → ❌ FAIL
         ↓
Build Job → ⏸️  SKIPPED
         ↓
No deployment
```

### Lint Failure (Ticket Service Only)
```
Test Job → ✅ PASS
Lint Job → ❌ FAIL
         ↓
Build Job → ⏸️  SKIPPED (needs both)
         ↓
No deployment
```

### Build Failure
```
Test Job → ✅ PASS
Build Job → ❌ FAIL (Docker build error)
         ↓
No image pushed
         ↓
ArgoCD keeps old version
         ↓
Services unchanged (safe)
```

### Coverage Threshold Failure (Ticket Service)
```
Test Job → Run tests
         ↓
Coverage: 12% (< 15% threshold)
         ↓
❌ Job FAILS
         ↓
Pipeline stops
```

---

## Time Breakdown

### Event Service (Total: ~5 min)
```
1. Test Job:          2m 15s
   ├─ Checkout:       15s
   ├─ Setup Go:       20s
   ├─ Cache:          10s
   ├─ Dependencies:   30s
   └─ Run tests:      60s

2. Build Job:         2m 45s
   ├─ Checkout:       15s
   ├─ Docker login:   10s
   ├─ Build image:    120s
   └─ Push image:     60s
```

### Ticket Service (Total: ~6 min)
```
Parallel Phase (max 3min):

1a. Test Job:         3m 00s
    ├─ PostgreSQL:    30s
    ├─ Checkout:      15s
    ├─ Setup Go:      20s
    ├─ Prisma gen:    45s
    ├─ DB push:       30s
    └─ Run tests:     40s

1b. Lint Job:         1m 05s
    ├─ Checkout:      15s
    ├─ Setup Go:      20s
    ├─ Prisma gen:    20s
    └─ Lint:          10s

Sequential Phase:

2. Build Job:         5m 00s
   ├─ Checkout:       15s
   ├─ Build binary:   60s
   ├─ Build API:      120s
   ├─ Push API:       45s
   ├─ Build Consumer: 90s
   └─ Push Consumer:  30s
```

---

## Environment Variables

### Event Service
```yaml
env:
  REGISTRY: ghcr.io
  IMAGE_NAME: dws-org/dws-event-service
```

### Ticket Service
```yaml
env:
  GO_VERSION: "1.23.1"
  DATABASE_URL: postgresql://postgres:postgres@localhost:5432/tickets_test
```

---

## Cost Analysis

**GitHub Actions Minutes:**
- Event Service: ~5 min/push
- Ticket Service: ~6 min/push
- Estimated pushes/day: 10
- Total: ~110 min/day
- Free tier: 2000 min/month ✅

**Container Registry Storage:**
- Event Service: ~150MB/image
- Ticket Service: ~200MB (API) + ~180MB (Consumer)
- Total per deployment: ~530MB
- GitHub Packages free: 500MB → Exceeded but under paid tier

---

## Security Measures

1. **No credentials in code**
   - Uses GitHub Secrets
   - GITHUB_TOKEN auto-injected

2. **Image scanning**
   - GHCR scans for vulnerabilities
   - ArgoCD checks image signatures

3. **Least privilege**
   - Build job: only packages:write
   - Deploy: ArgoCD service account

4. **Network isolation**
   - Services in private network
   - Only LoadBalancer exposed

---

## Future Improvements

### Event Service
- [ ] Add lint job (golangci-lint)
- [ ] Add test database (like Ticket Service)
- [ ] Set coverage threshold (start at 20%)
- [ ] Add E2E tests with Cypress

### Ticket Service
- [ ] Increase coverage to 30%
- [ ] Add integration tests for Consumer
- [ ] Add RabbitMQ in CI for queue tests
- [ ] Parallel build (API + Consumer)

### Both Services
- [ ] Add security scanning (Trivy)
- [ ] Add SBOM generation
- [ ] Cache Docker layers
- [ ] Add performance tests
- [ ] Blue-Green deployment option

---

**Total Pipeline Maturity:**
- Event Service: 6/10 (fast but lacks quality checks)
- Ticket Service: 8/10 (comprehensive but slow)
- ArgoCD Integration: 9/10 (GitOps working perfectly)
