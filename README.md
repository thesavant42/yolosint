# Forked from and Inspired by `oci.dag.dev`

This is a web server for exploring the contents of an OCI registry.

## Running it

Some things will probably break if the environment is different.

For local testing, I usually:

```
CACHE_DIR=/tmp/oci go run ./cmd/oci -v
```

On Cloud Run, I set `CACHE_BUCKET` to a GCS bucket in the same region as the service.

If you want private GCP images to work via oauth, you need to set `CLIENT_ID`, `CLIENT_SECRET`, and `REDIRECT_URL` to the correct values.

If you want to use ambient creds, set `AUTH=keychain`.

See also [`apk.dag.dev`](./cmd/apk/README.md);
