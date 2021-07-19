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
</p>

# Install
```
GO111MODULE=on go get -u github.com/xm1k3/cent
```
after installation run `cent init` to initialize cent with the configuration files you find [here](https://github.com/xm1k3/cent/blob/main/.cent.yaml) 


# Supported commands

| Command | Description            |
| ------- | ---------------------- |
| init    | Cent init project      |
| update  | Update your repository |

# Root flags

```
Flags:
      --config string   config file (default is $HOME/.cent.yaml)
  -C, --console         Print console output
  -k, --keepfolders     Keep folders (by default it only saves yaml files)
  -n, --name string     Name of the main folder
  -p, --path string     Path to save the templates
```

# Update flags
```
Flags:
  -d, --directories   Remove unnecessary folders from updated $HOME/.cent.yaml
  -f, --files         Remove unnecessary files from updated $HOME/.cent.yaml
```

# Init flags
```
Flags:
  -f, --field string   field to retrieve, comma separated
  -o, --overwrite      If the cent file exists overwrite it
  -u, --url string     Config folder url
```

# Usage

```
▶ cent -h
▶ cent update -h
```
Example
```
▶ cent -p {YOUR PATH} -n cent-nuclei-templates -k
▶ cent update -p {YOUR PATH} -d -f
```

Once cent has been configured correctly you can perform a scan with Nuclei.

Example
```
▶ nuclei -u https://example.com -t ./cent-nuclei-templates -tags cve
▶ nuclei -l urls.txt -t ./cent-nuclei-templates -tags cve
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
  - LICENCE

# Add github urls
community-templates:
  - https://github.com/geeknik/the-nuclei-templates
  - https://github.com/pikpikcu/nuclei-templates
  - https://github.com/panch0r3d/nuclei-templates
  - https://github.com/ARPSyndicate/kenzer-templates
  - https://github.com/medbsq/ncl
  - https://github.com/foulenzer/foulenzer-templates
  - https://github.com/im403/nuclei-temp
  - https://github.com/System00-Security/backflow
  - https://github.com/geeknik/nuclei-templates-1
  - https://github.com/esetal/nuclei-bb-templates
  - https://github.com/notnotnotveg/nuclei-custom-templates
  - https://github.com/clarkvoss/Nuclei-Templates
  - https://github.com/z3bd/nuclei-templates
  - https://github.com/joanbono/nuclei-templates
  - https://github.com/peanuth8r/Nuclei_Templates
  - https://github.com/thebrnwal/Content-Injection-Nuclei-Script
  - https://github.com/ree4pwn/my-nuclei-templates
  - https://github.com/optiv/mobile-nuclei-templates
  - https://github.com/obreinx/nuceli-templates
  - https://github.com/randomstr1ng/nuclei-sap-templates
  - https://github.com/CharanRayudu/Custom-Nuclei-Templates
  - https://github.com/zinminphyo0/KozinTemplates
  - https://github.com/n1f2c3/mytemplates
  - https://github.com/kabilan1290/templates
  - https://github.com/smaranchand/nuclei-templates
  - https://github.com/Saimonkabir/Nuclei-Templates
  - https://github.com/yavolo/nuclei-templates
  - https://github.com/sadnansakin/my-nuclei-templates
  - https://github.com/5cr1pt/templates
  - https://github.com/rahulkadavil/nuclei-templates
  - https://github.com/Nithissh0708/Custom-Nuclei-Templates
  - https://github.com/shifa123/detections
  - https://github.com/shifa123/mytemplates
  - https://github.com/daffainfo/my-nuclei-templates
```

# Want to help?
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/purple_img.png)](https://www.buymeacoffee.com/xm1k3)

## Credits
- [Alra3ees - Emad Shanab](https://twitter.com/Alra3ees)
  - [Nuclei-Templates-Collection](https://github.com/emadshanab/Nuclei-Templates-Collection)
- [sec715](https://twitter.com/sec715)
- [geeknik](https://twitter.com/geeknik)
- [Nuclei](https://twitter.com/pdnuclei)
- [Project Discovery](https://twitter.com/pdiscoveryio)
- [SYSTEM00 SECURITY](https://github.com/System00-Security)
- [clarkvoss](https://github.com/clarkvoss)
- [notnotnotveg](https://github.com/notnotnotveg)

# License
Cent is distributed under Apache-2.0 License