![Cent](./static/img/Cent_banner.png)

Community edition nuclei templates, a simple tool that allows you to organize all the Nuclei templates offered by the community in one place.

<p align="center">
<br>
<a href="https://github.com/xm1k3/cent/issues"><img src="https://img.shields.io/badge/contributions-welcome-success.svg?style=flat"></a>
<img alt="Apache license badge" src="https://img.shields.io/badge/license-Apache-success">
<a href="https://github.com/xm1k3/cent/releases"><img src="https://img.shields.io/github/release/xm1k3/cent"></a>
<br>
<a href="https://github.com/xm1k3/cent/stargazers"><img src="https://img.shields.io/github/stars/xm1k3/cent.svg?style=social&label=Stars"></a>
<a href="https://twitter.com/xm1k3_"><img src="https://img.shields.io/twitter/follow/xm1k3_.svg?logo=twitter"></a>
<br>
<br>
<br>
<a href="https://www.buymeacoffee.com/xm1k3"><img src="https://www.buymeacoffee.com/assets/img/custom_images/purple_img.png"></a>
</p>

# Install
```
go install -v github.com/xm1k3/cent@latest
```

Or [download from releases](https://github.com/xm1k3/cent/releases)
<br><br>

after installation run `cent init` to initialize cent with the configuration files you find [here](https://github.com/xm1k3/cent/blob/main/.cent.yaml) 


# Supported commands

| Command | Description            |
| ------- | ---------------------- |
| check | Check if templates repo are still available |
| init    | Cent init configuration file      |
| summary | Print summary table | 
| update  | Update your repository |
| validate | Validate templates, if the template is invalid it is deleted from the folder |
| version  | Print cent version |

# Root flags

```
Flags:
      --config string   config file (default is $HOME/.cent.yaml)
  -C, --console         Print console output
  -p, --path string     Root path to save the templates (default "cent-nuclei-templates")
  -t, --threads int     Number of threads to use when cloning repositories (default 10)
```


# Usage

```
cent -h
cent check -h
cent init -h
cent update -h
cent summary -h
cent validate -h
cent version
```
Example:

Clone and insert all the community templates into the `cent-nuclei-templates` folder 
```
cent -p cent-nuclei-templates
```
![cent](./static/img/cent-v1.0.png)

If you have updated the `cent.yaml` file by adding new folders
```yaml
exclude-dirs:
  - ...
  - dns
  - ...
```
just do:
```
cent update -p cent-nuclei-templates -d
```
and `cent` will automatically delete all `dns` folder present in `cent-nuclei-templates` without cloning all the github repos.

![cent update](./static/img/cent-update.png)

Same thing with `exclude-files`
```
cent update -p cent-nuclei-templates -f
```
---
Once cent has been configured correctly you can perform a scan with Nuclei.

Example
```
nuclei -u https://example.com -t ./cent-nuclei-templates -tags cve
nuclei -l urls.txt -t ./cent-nuclei-templates -tags cve
```
See [here](https://nuclei.projectdiscovery.io/nuclei/get-started/#running-nuclei) for more documentation about Nuclei


# Config
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
  - LICENSE

# Add github urls
community-templates:
  - https://github.com/projectdiscovery/nuclei-templates
  ...
  ...

```

## Credits
- [hakluke](https://twitter.com/hakluke)
- [Nuclei](https://twitter.com/pdnuclei)
- [Project Discovery](https://twitter.com/pdiscoveryio)
- [sec715](https://twitter.com/sec715)
- [geeknik](https://twitter.com/geeknik)
- [SYSTEM00 SECURITY](https://github.com/System00-Security)
- [clarkvoss](https://github.com/clarkvoss)
- [notnotnotveg](https://github.com/notnotnotveg)
- [Alra3ees - Emad Shanab](https://twitter.com/Alra3ees)
- [Nuclei-Templates-Collection](https://github.com/emadshanab/Nuclei-Templates-Collection)


# Disclaimer

Disclaimer: The developer of this tool is not responsible for how the community uses the open source templates collected within it. These templates have not been validated by Project Discovery and are provided as-is.

# License
Cent is distributed under Apache-2.0 License
