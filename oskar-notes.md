### Was muss man machen?

1. Docker Image lokal bauen  
   ```bash
   docker build -t dws-event-service:latest .
   ```
2. Image für Docker Hub taggen & pushen  
   ```bash
   # Beispiel: Docker Hub User "oskar"
   docker tag dws-event-service:latest oskar/dws-event-service:latest
   docker push oskar/dws-event-service:latest
   ```
3. Image für GitHub Container Registry (GHCR) taggen & pushen  
   ```bash
   # vorher einmal bei ghcr einloggen:
   # echo "$GHCR_TOKEN" | docker login ghcr.io -u <github-username> --password-stdin

   docker tag dws-event-service:latest ghcr.io/<github-username>/dws-event-service:latest
   docker push ghcr.io/<github-username>/dws-event-service:latest
   ```
4. `values.yaml` im Helm-Chart anpassen (Image + ggf. Tag)  
   ```yaml
   image:
     repository: ghcr.io/<github-username>/dws-event-service
     tag: "latest"
   ```
5. Mit Helm auf Kubernetes deployen / updaten  
   ```bash
   # kube-context + namespace ggf. anpassen
   kubectl config use-context <dein-k8s-context>
   kubectl create namespace dws --dry-run=client -o yaml | kubectl apply -f -

   helm upgrade --install dws-event-service ./helm/dws-event-service -n dws
   ```

## Build (Kurz-Checkliste)
1. `docker build -t dws-event-service:latest .`  
2. Tags für Docker Hub & GHCR setzen (`docker tag ...`)  
3. `docker push` für beide Registries  
4. `helm/dws-event-service/values.yaml` → `image.repository` + `image.tag` anpassen  
5. `helm upgrade --install ...` ausführen
