# Refactor Image

I want to tidy up the image view but I do not want to break anything.

## Summary 

I want the [current view](/plans/current.png) to look like [this](/plans/mockup.png)  instead.

### Key Differences:

5 columns -> **4 Columns:**

#### Columns fixed
- 1. Index Layer `#`
- 2. `Size`
- 3. `sha256:digest`
- 4. Text link for layer downloads `[X]`
**Critical insight: Remove redundant index column, making 5 columns -> 4 columns**

### But then, theese changes also:

- 5. `medaType...` is shortened to `Manifest v2`
- 6. `[LIST VIEWS]` label extends past the first column, over the second column, and the label **does NOT WRAP**
- 7. "`CONFIG`" 
    - ~~label is capitalized~~ 
    - Aligned **vertically** with the **row numbers** 
- 8. '`[COMBINED LAYERS]` view moved to the roqw with the shortened `manifest` tag, and before the moved `referrers` component.

## Acceptance Criteria:
1. Must look like the screenshot
2. `COMBINED LAYERS`, `CONFIG`, `LIST VIEW`, and `LAYERS VIEW`, `SAVE LAYERS` links **must cotinue to work.**


### Validation links:

1. [Overview link](http://192.168.1.37:8042/?image=pullu84%2Fballys-roadmap-assistant%40sha256%3A72a32dab46f13f0d277ec98295d6c5371ce895f13718c724f065d6d5df95d6fc)

2. [Referral](http://192.168.1.37:8042/?referrers=pullu84/ballys-roadmap-assistant@sha256:72a32dab46f13f0d277ec98295d6c5371ce895f13718c724f065d6d5df95d6fc)

3. [CONFIG](http://192.168.1.37:8042/?blob=pullu84/ballys-roadmap-assistant@sha256:9426118caaa40ecd130edf679c1744400a70a08cd6aaf987f78b44c20184519e&mt=application%2Fvnd.oci.image.config.v1%2Bjson&size=9083)

4. [list views](http://192.168.1.37:8042/size/pullu84/ballys-roadmap-assistant@sha256:12431f47c511cdc7010d4c0c0bc4de33fd45b5167f14e6724e324c5982d6e04f?mt=application%2Fvnd.oci.image.layer.v1.tar%2Bgzip&size=1095717)

5. [layer9](http://192.168.1.37:8042/fs/pullu84/ballys-roadmap-assistant@sha256:12431f47c511cdc7010d4c0c0bc4de33fd45b5167f14e6724e324c5982d6e04f/?mt=application%2Fvnd.oci.image.layer.v1.tar%2Bgzip&size=1095717)

6. [save layer 9](http://192.168.1.37:8042/download/pullu84/ballys-roadmap-assistant@sha256:12431f47c511cdc7010d4c0c0bc4de33fd45b5167f14e6724e324c5982d6e04f?filename=pullu84-ballys-roadmap-assistant-9.tar.gzip)

7. [COMBINED LAYERS](http://192.168.1.37:8042/layers/pullu84/ballys-roadmap-assistant@sha256:72a32dab46f13f0d277ec98295d6c5371ce895f13718c724f065d6d5df95d6fc/) 