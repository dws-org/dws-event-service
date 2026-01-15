# CI/CD Pipeline Diagrams

This folder contains CI/CD pipeline visualizations for the DWS platform.

## Files

1. **c4-cicd-pipeline.puml** - C4 Deployment diagram (high-level flow)
2. **c4-cicd-linear.puml** - Standard sequence diagram (detailed linear flow)
3. **CICD_PIPELINE_FLOW.md** - Comprehensive documentation with ASCII diagrams

## How to View PlantUML Diagrams

### Online (Easiest)
**C4 Deployment Diagram:**
1. Go to http://www.plantuml.com/plantuml/uml/
2. Copy content from `c4-cicd-pipeline.puml`
3. Paste and click "Submit"

**Sequence Diagram:**
1. Go to http://www.plantuml.com/plantuml/uml/
2. Copy content from `c4-cicd-linear.puml`
3. Paste and click "Submit"

### VS Code
1. Install "PlantUML" extension (by jebbs)
2. Open `.puml` file
3. Press `Alt+D` to preview (or right-click → "Preview Current Diagram")

### Command Line
```bash
# Install PlantUML (requires Java)
sudo apt install plantuml  # Ubuntu/Debian
brew install plantuml      # macOS

# Generate images
plantuml c4-cicd-pipeline.puml
plantuml c4-cicd-linear.puml

# Output: PNG images in same directory
```

## Quick Overview

```
Developer → GitHub → Actions → GHCR → ArgoCD → Kubernetes
            (code)   (CI/CD)   (images)  (GitOps) (deploy)
```

**Event Service Pipeline:** Test → Build (2 jobs, ~3 min)
- ✅ Prisma client generation added
- ✅ 20.5% test coverage
- ✅ Codecov integration

**Ticket Service Pipeline:** Test + Lint → Build (3 jobs, ~6 min)
- ✅ PostgreSQL test database
- ✅ golangci-lint checks
- ✅ Dual image build (API + Consumer)

See CICD_PIPELINE_FLOW.md for detailed breakdown.

## Diagram Details

### c4-cicd-pipeline.puml (C4 Deployment)
- **Type:** C4 Level 3 - Component/Deployment diagram
- **Library:** C4-PlantUML (C4_Deployment.puml)
- **Layout:** Left-to-right flow
- **Shows:** 
  - System boundaries (GitHub, CI/CD, Kubernetes, Monitoring)
  - Service flows for both Event and Ticket services
  - Container relationships
  - Deployment pipeline stages

### c4-cicd-linear.puml (Sequence)
- **Type:** Standard PlantUML sequence diagram
- **Style:** UML sequence with participants
- **Shows:**
  - Chronological step-by-step flow
  - Parallel job execution
  - Time-based interactions
  - Test, lint, build phases with actual coverage numbers

### CICD_PIPELINE_FLOW.md (Documentation)
- **Type:** Markdown with ASCII diagrams
- **Contains:**
  - 3 detailed flow diagrams
  - Time breakdown per job
  - Failure scenarios
  - Cost analysis
  - Comparison tables

## Testing the Diagrams

Both diagrams have been tested and work correctly:
- ✅ C4 deployment diagram renders successfully
- ✅ Sequence diagram uses standard PlantUML syntax (no C4_Sequence dependency)
- ✅ All relationships and notes display properly

Latest CI/CD pipeline run: ✅ Successful (2m32s)
