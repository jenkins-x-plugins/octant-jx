### Linux

```shell   
mkdir -p  ~/.config/octant/plugins
curl -L https://github.com/jenkins-x/octant-jx/releases/download/v{{.Version}}/octant-jx-linux-amd64.tar.gz | tar xzv 
mv octant-* ~/.config/octant/plugins
```

### macOS

```shell
mkdir -p  ~/.config/octant/plugins
curl -L  https://github.com/jenkins-x/octant-jx/releases/download/v{{.Version}}/octant-jx-darwin-amd64.tar.gz | tar xzv
mv octant-* ~/.config/octant/plugins
```

