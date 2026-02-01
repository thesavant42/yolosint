# Docker Dorker - OCI ICU 

Forked from and Inspired by `oci.dag.dev`

## Docker Dorker
This is a web server for exploring the contents of an OCI registry, FS style output.

## Running it
`rm ./oci; git pull; go build /cmd/oci ; ./oci -v;`

Opens a listener on port localhost:8080 but I forward this sto port 8042 in Docker. That aligns with the url on line 16 of [DockerHubOCI Explorer-1.1.user.js](/DockerHubOCI Explorer-1.1.user.js)