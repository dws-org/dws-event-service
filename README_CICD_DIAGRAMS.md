# CI/CD Pipeline Diagrams

This folder contains CI/CD pipeline visualizations for the DWS platform.

## Files

1. **c4-cicd-pipeline.puml** - C4 Deployment diagram (high-level flow)
2. **c4-cicd-linear.puml** - C4 Sequence diagram (detailed linear flow)
3. **CICD_PIPELINE_FLOW.md** - Comprehensive documentation with ASCII diagrams

## How to View PlantUML Diagrams

### Online
1. Go to http://www.plantuml.com/plantuml/uml/
2. Copy & paste the `.puml` file content
3. Click "Submit"

### VS Code
1. Install "PlantUML" extension
2. Open `.puml` file
3. Press `Alt+D` to preview

### Command Line
```bash
plantuml c4-cicd-pipeline.puml
# Generates: c4-cicd-pipeline.png
```

## Quick Overview

```
Developer → GitHub → Actions → GHCR → ArgoCD → Kubernetes
            (code)   (CI/CD)   (images)  (GitOps) (deploy)
```

**Event Service Pipeline:** Test → Build (2 jobs, ~5 min)
**Ticket Service Pipeline:** Test + Lint → Build (3 jobs, ~6 min)

See CICD_PIPELINE_FLOW.md for detailed breakdown.
