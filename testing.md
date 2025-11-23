kubectl run test-events-api --image=curlimages/curl --rm -i --restart=Never -- \
  curl -s http://dws-event-service.dws-event-service.svc.cluster.local:6906/api/v1/events

  