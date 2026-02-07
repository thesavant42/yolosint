# Docker Dorker - YOLOSINT Container Module 

<img src="docs/128-cropped-scaled-banner-001.png" style="display: block; margin: 0 auto;" alt="logbanner">

Forked from and Inspired by [https://github.com/jonjohnsonjr/dagdotdev](https://github.com/jonjohnsonjr/dagdotdev)

## Docker Dorker
This is a web server for exploring the contents of an OCI registry, FS style output.

<img src="docs/screenshots/rsapriv.png" style="display: block; margin: 0 auto;" alt="rsa smile">

---

## Quick Start (Docker)

```bash
git pull
docker compose up --build
```

Service runs at `http://localhost:8042`


**Now with verbose logging!**

---

## Running Locally (without Docker)

```bash
git pull
rm ./oci
go build ./cmd/oci
./oci -v
```

Opens a listener on port *localhost:8080*, **but I forward this to port 8042** in Docker.
    - That aligns with the url on line 16 of [DockerHubOCI yolosint.user.js](/DockerHubOCI Explorer-1.1.user.js)

## Cached Files
Stored in the [./cache/](./cache/) folder.

---

## User Script
- `yolosint.user.js` is a user script written for Tampermonkey, ymmv.
- Adds an overlay when viewing docker hub to quickly view a container image in docker dorker. 
- To use after installing, simply browse an image in Docker Hub and click the "Tags" tab to view details about recent tags for that Repository.
- The overlay will be displayed nexted to the Tag ID.

---

## Workflow

### 1. Landing Page

- Main View of API Landing Page
- Mine is installed on my NAS via Docker


<img src="docs/screenshots/main.png" style="display: block; margin: 0 auto;" alt="main">


### 2. Searching Docker Hub

Searching for "contractor" using the landing page helper. I can search directly from the landing page for the service, or I can search directly from docker hub.


<img src="docs/screenshots/search-results.png" style="display: block; margin: 0 auto;" alt="search results">



Results appear as per usual.

### 3. View Tags for interesting container repositories

Once I locate a container repository I am interested in, I need to view the Tags.

<img src="docs/screenshots/tags-results.png" style="display: block; margin: 0 auto;" alt="tags results">


### 4. User Script provides Overlay to view Layer in the Docker Dorker module.
The Tampermonkey user script overlays a help "YOLOSINT" banner to jump into analsys view.



<img src="docs/screenshots/extension.png" style="display: block; margin: 0 auto;" alt="tags extension">


### 5. Overview of Tag Layers


<img src="docs\screenshots\sample-tagview.png" style="display: block; margin: 0 auto;" alt="sample tag view">

Multiple architectures and image types are available in this repository.

### 6. View Tag Cofnig / Build Steps Info

The Image Config manifest contains steps used to build the container. This is often helpful; environment variables being set with credentials, WORKINGDIR to know which filesystem paths are of interest, the ENTRYPOINT, the exposed PORTS.



<img src="docs/screenshots/config-example.png" style="display: block; margin: 0 auto;" alt="config example">


### 7. View an overview of the Filesystem Layers

### 7. View an overview of the Filesystem Layers

The container image itself is created via an Overlay filesystem, which consists of multiple filesystem image "layers", stored as .tar.gz files and addressed by their SHA256 Digest indentifiers.

<img src="docs/screenshots/example-fsview.png" style="display: block; margin: 0 auto;" alt="sample fs view">

### 8. View the details of a single Layer at a time

Indexing the filesyste of the container with YOLOSINT's Docker Dorker module take a fraction of a second. I can view layers ine at a time, like this one.

<img src="docs/screenshots/layerview.png" style="display: block; margin: 0 auto;" alt="singleLayerView">

Or I can view the Merged Filesystem view, which merges them all together. This view allows me to click through the filesystem, presented as a simulated "ls -la" terminal command output.

### 9. Or view all layers, merged FS

Merged FS View

<img src="docs/screenshots/mergedview.png" style="display: block; margin: 0 auto;" alt="mergedfsview">

### 10. View Layers in "Size" view, which lists the files and their details.

If that's not preferred, viewing the contents of a layer as a detailed list is also possible.

<img src="docs/screenshots/size-view.png" style="display: block; margin: 0 auto;" alt="sizeview">

### 11. Binary File review

Let's suppose we want to analyze the `su` binary on the filesystem.

<img src="docs/screenshots/suview2.png" style="display: block; margin: 0 auto;" alt="binaryview">

### 12. Binary File Review - Details

Clicking the file provides a basic HEX "xxd" view of the file.

<img src="docs/screenshots/binaryview.png" style="display: block; margin: 0 auto;" alt="binaryfiledetails">

### 13. Binary File Review - Details (Containued)

Clicking the ELF hyperlink launches detailed analysis view.

<img src="docs/screenshots/binarydetails.png" style="display: block; margin: 0 auto;" alt="binarydetailscontd">

### 14. Binary Analysis View

Without any docker container saving, we're able to determine the contents of the file, how it was built, and more. 

<img src="docs/screenshots/bindetails-golang.png" style="display: block; margin: 0 auto;" alt="binaryanalysis">

### 15. Save the Binary

<img src="docs/screenshots/savesu.png" style="display: block; margin: 0 auto;" alt="savebinary">

### 16. Plaintext Oopsies

Some real-world examples, redacted to protect the accidental. **This example has a basic AUTH header to pull from a private GitHuvb repository.**

<img src="docs/screenshots/oopsies.png" style="display: block; margin: 0 auto;" alt="oopsies1">

#### Ooopsies!

**Oopsies, Firebase authentication and Google IAM credentials.**

<img src="docs/screenshots/oopsies-creds.png" style="display: block; margin: 0 auto;" alt="oopsies2">

### 17. Save as Preview

Viewing the content of "saved" files in the browser can be helpful, such as when rendering images embedded in a container.

<img src="docs/screenshots/save-magic.png" style="display: block; margin: 0 auto;" alt="save-magic">

I can preview the icon by saving it.

<img src="docs/screenshots/savemagic-icon.png" style="display: block; margin: 0 auto;" alt="previewicon">