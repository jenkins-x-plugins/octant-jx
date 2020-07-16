# Jenkins X Octant Plugins

This repository contains plugins for [Octant](https://octant.dev/) for working with [Jenkins X](https://jenkins-x.io/)

## Install

### Prerequisites

Octant should first be installed and added to your `$PATH`, get the latest release [here](https://github.com/vmware-tanzu/octant/releases)

### Install the plugin

Octant checks for extra plugins that live in `~/.config/octant/plugins`.  You can download the released plugin binary [here](https://github.com/jenkins-x/octant-jx/releases/) or build it locally in the next section. 

### Development of the plugin

To build the plugin use:

```
make octant
```

which will build the plugin, install it into `~/.config/octant/plugins` and then startup octant against the current k8s cluster.

You can run `make tail` in another terminal to watch the console log of `octant-jx` if you are developing a plugin.

### Using the Developer + Operator plugins

You should see on the left nav bar the Jenkins X Developer plugin appears near the bottom (2nd to last icons).
 
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
