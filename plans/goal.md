# Goal  Simpify Search Boxes from 3 to 1

Problem:

Search is a confusing mess on the app mainpage as seen in this screenshot:

![threes search boxes](/docs/screenshots/main.png)


There are threes search input boxes, but only one can be used at any give time. 

## Proposed Solution:

### Selection Dropdown or Radio

Create a Single Text Input, 3 routes:
    1. Image
        Requests to `/?image=namespace/repo:tag`
    2. Repo
        Requests to `/?repo=repo`
    3. Docker Hub
        POSTs to `https://hub.docker.com/search?q=`
    ~~4. (TODO: Cache History Search)~~
        ~~1. Search the sqlite database for previously observed file names~~

#### Considerations:

- Must still allow DIRECT posts and gets to continue to work 
- example text will need to be refactored to provider hovertext help or to render it when the option is selected

- Default to optionn 1, the Image route (currently the top most text input in the screenshot)
    - `<input size="100" type="text" name="image" value="ubuntu:latest">`