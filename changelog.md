### Linux

```shell   
mkdir -p  ~/.config/octant/plugins
curl -L https://github.com/jenkins-x/octant-jx/releases/download/v0.0.44/octant-jx-linux-amd64.tar.gz | tar xzv 
mv octant-* ~/.config/octant/plugins
```

### macOS

```shell
mkdir -p  ~/.config/octant/plugins
curl -L  https://github.com/jenkins-x/octant-jx/releases/download/v0.0.44/octant-jx-darwin-amd64.tar.gz | tar xzv
mv octant-* ~/.config/octant/plugins
```

## Changes

### Bug Fixes

* pipeline (James Strachan)
* lets visualise boot jobs nicer (James Strachan)

### Chores

* fix updatebot (James Strachan)
