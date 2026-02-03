# Docker Hub Buil+ Push w/ SBOM Attestation

```bash
docker buildx build --sbom=true --provenance=mode=max -t savant42/yolosint:latest --push .
```

