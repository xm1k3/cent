# cent
Community edition nuclei templates, a simple tool that allows you to organize all the Nuclei templates offered by the community in one place

# Supported commands

| Command            | Description                  |
|--------------------|------------------------------|
| path                 | Path to save the templates   |
| name               | Name of the main folder   |
| keepfolders              | Keep folders (by default it only saves yaml files)        |
| console             | Print console output                   |

# Install
```
go get -u github.com/xm1k3/cent
```

# Usage

```
cent -h
go run main.go -h
```
Example
```
▶ cent -p ~/ -n community-nuclei-templates -k
▶ cent -p {YOUR PATH} -n community-nuclei-templates
```

### Root flags
```
Flags:
      --config string   config file (default is $HOME/.cent.yaml)
  -C, --console         Print console output
  -h, --help            help for cent
  -k, --keepfolders     Keep folders (by default it only saves yaml files)
  -n, --name string     Name of the main folder
  -p, --path string     Path to save the templates
```

### Config
You need to configure `cent` parameters in `$HOME/.cent.yaml`
```yaml
exclude-dirs:
  - .git

exclude-files:
  - README.md
  - .gitignore
  - .pre-commit-config.yaml
  - LICENCE

# Add github urls
community-templates:
  - https://github.com/geeknik/the-nuclei-templates.git
  - https://github.com/pikpikcu/nuclei-templates.git
  - https://github.com/panch0r3d/nuclei-templates.git
  - https://github.com/ARPSyndicate/kenzer-templates.git
  - https://github.com/medbsq/ncl.git
  - https://github.com/foulenzer/foulenzer-templates.git

```

## Credits
- [Alra3ees](https://twitter.com/Alra3ees)
- [sec715](https://twitter.com/sec715)
- [geeknik](https://twitter.com/geeknik)
- [Nuclei](https://twitter.com/pdnuclei)
- [Project Discovery](https://twitter.com/pdiscoveryio)