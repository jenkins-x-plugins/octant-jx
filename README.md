# Jenkins X Octant Plugins

This repository contains plugins for [Octant](https://octant.dev/) for working with [Jenkins X](https://jenkins-x.io/)

## Why Octant

[Octant](https://github.com/vmware-tanzu/octant) is the strategic UI for working with Jenkin X because its:

* open source and very easy to extend with plugins
* lets you visualise and work with all kubernetes and custom resources across multiple clusters
* thanks to `octant-jx` has nice integration with Jenkins X components like apps, environments, pipelines, repositories etc.
* over time will add management UI capabilities for installing, upgrading and administering Jenkins  

## Demo

Here is a [demo video showing octant in action with Jenkins X](https://www.youtube.com/watch?v=2LCPHi0BnUg&feature=youtu.be):

<iframe width="1292" height="654" src="https://www.youtube.com/embed/2LCPHi0BnUg" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
  
We also [presented octant-jx](https://www.youtube.com/watch?v=Njl247hjRuU&t=2027s) at the [octant office hours this week](https://octant.dev/community/).

## Features

Longer term we're planning on making most of the developer and operations feaures of Jenkins X available through the UI via [octant-jx](https://github.com/jenkins-x/octant-jx) but already you can:

* view applications, environments, pipelines, repositories
* for a pipeline quickly navigate to:
  * its Pod, Log, Pull Request or Preview Environment
  * for each step you can view the step detail or log of the step
* see the various jobs and pipelines used to operate Jenkins X itself
* over time will add management UI capabilities for installing, upgrading and administering Jenkins  

## Install

### Prerequisites

Octant should first be installed and added to your `$PATH`.

Get the latest release from [vmware-tanzu/octant](https://github.com/vmware-tanzu/octant/releases)

### Install the plugin

Octant checks for extra plugins that live in `~/.config/octant/plugins`.  

You can download the released plugin binaries [here](https://github.com/jenkins-x/octant-jx/releases/) and move the `octant-*` binaries to` ~/.config/octant/plugins`

## Running

Run the UI via:

```bash 
octant --browser-path="/#/jx/pipelines"
```

You should see on the left nav bar the Jenkins X Developer + Ops plugins appear near the bottom (2nd to last icons).
 
## Running multiple Octants

You may connect to different clusters in different shells and open an octant for each cluster via:

``` 
octant --listener-addr=localhost:0
```

If you want to open a specific view try:

``` 
octant --listener-addr=localhost:0  --browser-path="/#/jx/pipelines-recent"
```

An octant will start along with a new browser window.


## Developing octant-jx 

`octant-jx` is 100% go lang and has a pretty simple small code base - so we'd love contributions! 

It should be easy to add or improve the UI to handle most use cases.

### Building locally
 
To build the plugins use:

```
make octant
```

which will build the plugins, install then into `~/.config/octant/plugins` and then startup octant against the current k8s cluster.

You can run `make tail` in another terminal to watch the console log of `octant-jx` if you are developing a plugin.
 

## Adding a new view

Its super easy to add new views of any kubernetes or custom resource you like.

A good example to copy/paste if you want to add a new view is the `Repositories` view:

* [pkg/plugin/views/repositories_view.go](https://github.com/jenkins-x/octant-jx/blob/master/pkg/plugin/views/repositories_view.go) is the view itself then you need to:
  * add a handler in [pkg/plugin/router/handlers.go](https://github.com/jenkins-x/octant-jx/blob/master/pkg/plugin/router/handlers.go#L24) to use your new view file
  * add a link to your new view in [pkg/plugin/settings/options.go](https://github.com/jenkins-x/octant-jx/blob/master/pkg/plugin/settings/options.go#L40-L44)


### Resources for plugin developers 

* Check out the [Plugin Documentation](https://octant.dev/docs/master/plugins/)
* Views can use any of the widgets in the `/pkg/view/components` package of octant - check the [reference docs](https://octant.dev/docs/master/plugins/reference/).
* Here's another [example plugin](https://github.com/vmware-tanzu/octant/blob/master/cmd/octant-sample-plugin/main.go#L27) which enriches existing views
