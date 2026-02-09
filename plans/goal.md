# Print path when browsing container

## Problem statement:

There's no quick way to know which layer a user is browsing at a given time. 

### Solution proposal:

I propose displaying the current layer number, before the path in the UI. 


#### Current View: 

![/plans/current.png](/plans/current.png)

#### Proposed View Mockup:

![plans/mockup.png](/plans/mockup.png)

- Shows the layer, above the path
- includes a download link to save the layer. This is the same link that's rendered next to the layer in the /image? root view.
    - [/internal/explore/assets/gis--layer-download.png](/internal/explore/assets/gis--layer-download.png) this is the icon to include, 
    - I could not include it in the mockup, but it should go before the text of the download link.
- The Merged view would similarly include a label indicating the layers are merged.
- Where the download link would normally be per layer, the merged view can draw a `--` for merged view
- 
![plans/mockup.png](/plans/mockup2.png)



### Task:

- Read relevant files:
    - [/internal/explore/explore.go](/internal/explore/explore.go)
    - [/internal/explore/render.go](/internal/explore/render.go)
    - [/internal/explore/templates.go](/internal/explore/templates.go)


### Acceptance Critera:
- layer is printed above path with a download link when when single layer view
- layer idicates merged view when viewing the merged view, does *not* provide "layer download" link in header