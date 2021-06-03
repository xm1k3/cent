# cent
Community edition nuclei templates, a simple tool that allows you to organize all the Nuclei templates offered by the community in one place

# Install
```
go get -u github.com/xm1k3/cent
```

# Supported commands

| Command            | Description                  |
|--------------------|------------------------------|
| update                 | Update your repository   |

# Root flags

```
Flags:
      --config string   config file (default is $HOME/.cent.yaml)
  -C, --console         Print console output
  -h, --help            help for cent
  -k, --keepfolders     Keep folders (by default it only saves yaml files)
  -n, --name string     Name of the main folder
  -p, --path string     Path to save the templates
```

# Update flags
```
Flags:
  -d, --directories   Remove unnecessary folders from updated $HOME/.cent.yaml
  -f, --files         Remove unnecessary files from updated $HOME/.cent.yaml
  -h, --help          help for update
```

# Usage

```
cent -h
cent update -h
```
Example
```
▶ cent -p {YOUR PATH} -n community-nuclei-templates -k
▶ cent update -p {YOUR PATH} -d -f
```


### Config
You need to configure `cent` parameters in `$HOME/.cent.yaml`
```yaml
# Directories to exclude
exclude-dirs:
  - .git

# Files to exclude
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
  - https://github.com/im403/nuclei-temp
  - https://github.com/System00-Security/backflow
  - https://github.com/geeknik/nuclei-templates-1
  - https://github.com/esetal/nuclei-bb-templates
  - https://github.com/notnotnotveg/nuclei-custom-templates
  - https://github.com/clarkvoss/Nuclei-Templates

```

## Want to help?
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/xm1k3)

## Credits
- [Alra3ees](https://twitter.com/Alra3ees)
- [sec715](https://twitter.com/sec715)
- [geeknik](https://twitter.com/geeknik)
- [Nuclei](https://twitter.com/pdnuclei)
- [Project Discovery](https://twitter.com/pdiscoveryio)
- [SYSTEM00 SECURITY](https://github.com/System00-Security)
- [clarkvoss](https://github.com/clarkvoss)
- [notnotnotveg](https://github.com/notnotnotveg)