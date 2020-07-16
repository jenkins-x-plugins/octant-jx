### Linux

```shell
curl -L https://github.com/jenkins-x/octant-jx/releases/download/v{{.Version}}/octant-jx-amd64.tar.gz | tar xzv 
mv octant-jx ~/.config/octant/plugins
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/octant-jx/releases/download/v{{.Version}}/octant-jx-darwin-amd64.tar.gz | tar xzv
mv octant-jx ~/.config/octant/plugins
```

