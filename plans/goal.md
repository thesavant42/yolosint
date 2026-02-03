# Task Goal: Enhance output in layers view to make the layers more obvious in their intent

## Goal Statement

Enhance output of /image route

### Curret State 

Current state of `/?image=kentico%2Fems%3Aweb-12.0.29`

![currentstate](/plans/current.png) 

### Desired End State

Desired End State Mockup:

![MOCKUP](/plans/MOCKUP.png)

### Key differences:

- Added a column for each table that includes the numeric identifier starting at `1`, to make it easier to cross reference individual layers
    - I'd like these layer identifiers to hyperlink to download the entire layer image tar.gz file, but renamed to be human readable:
        - Instead of `sha256-9038b92872bc268d5c975e84dd94e69848564b222ad116ee652c62e0c2f894b2.tar.gz`
        - It would be: `repo-tag-layer-integer`, so `repo-tag-layer-1.tar.gz`
            - **There is NOT currently a route to download the layer individually**
- `layers` (`/layers/NAMESPACE/REPO:TAG@sha256:DIGEST/`) has added text label to make it `Combined Layers View` to make its intent more obvious
- Config sha256 digest hyperlink moved up to be on the same row as the word `config`
- Added text labels to tables to make their purpose more obvious
- `docker pull` copy box added (see ["New Teature Request" below])

### Key Insight: 

- The blue text is not a design choice, it was used to make the new text ore obvious in my explanation to you

---


## New Feature Request:

 I'd also lilke to add a new feature:
Docker Hub style "docker pull namespace/repo:tag" copy box, which prefills the users clipboard.

#### Add docker pull copy box like Docker Hub:

- The Docker hub overview page for a repository offers a "copy" box that prepopulates the users clipboard with the proper docker pull tab to download the entire container. 

- Box without mouse over hover effect:

![copybox](/plans/copybox.png)

- Box WITH mouse over effect:

![copyboxhover](/plans/copyboxhover.png)

- This is the HTML from that page, as seen on pages such as `https://hub.docker.com/r/kentico/ems/tags`
 
```html
<button class="MuiButtonBase-root MuiButton-root MuiButton-contained MuiButton-containedPrimary MuiButton-sizeMedium MuiButton-containedSizeMedium MuiButton-colorPrimary MuiButton-root MuiButton-contained MuiButton-containedPrimary MuiButton-sizeMedium MuiButton-containedSizeMedium MuiButton-colorPrimary css-1v11hr0" tabindex="0" type="button" data-testid="copy-code" style="right: 4px;">Copy</button>
```