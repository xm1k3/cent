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
| version  | Print cent version |

# Root flags

```
Flags:
      --config string   config file (default is $HOME/.cent.yaml)
  -C, --console         Print console output
  -k, --keepfolders     Keep folders (by default it only saves yaml files)
  -p, --path string     Root path to save the templates (default "cent-nuclei-templates")
  -t, --threads int     Number of threads to use when cloning repositories (default 10)
```

# Update flags
This command helps you update your folder with templates by deleting unnecessary folders and files without having to do multiples git clones.
```
Flags:
  -d, --directories   If true remove unnecessary folders from updated $HOME/.cent.yaml
  -f, --files         If true remove unnecessary files from updated $HOME/.cent.yaml
  -p, --path string   Path to folder with nuclei templates
```

# Init flags
This command will automatically download [`.cent.yaml`](https://raw.githubusercontent.com/xm1k3/cent/main/.cent.yaml) from repo and copy it to `$HOME/.cent.yaml`
```
Flags:
  -h, --help         help for init
  -o, --overwrite    If the cent file exists overwrite it
  -u, --url string   Url from which you can download the configurations for .cent.yaml
```

# Check flags
```
Flags:
  -h, --help     help for check
  -r, --remove   Remove from .cent.yaml urls that are no longer accessible or are currently private
```

# Summary flags

```
Flags:
  -h, --help          help for summary
  -p, --path string   Root path to save the templates (default "cent-nuclei-templates")

output: 

+-------------------------------+-----------------+
| TEMPLATES TYPE                | TEMPLATES COUNT |
+-------------------------------+-----------------+
| CVE Templates                 |           21981 |
| Other Vulnerability Templates |           16779 |
| Total Templates               |           38760 |
+-------------------------------+-----------------+
```


# Usage

```
▶ cent -h
▶ cent init -h
▶ cent update -h
▶ cent version
```
Example:

Clone and insert all the community templates into the `cent-nuclei-templates` folder 
```
▶ cent -p cent-nuclei-templates -k
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
▶ cent update -p cent-nuclei-templates -d
```
and `cent` will automatically delete all `dns` folder present in `cent-nuclei-templates` without cloning all the github repos.

![cent update](./static/img/cent-update.png)

Same thing with `exclude-files`
```
▶ cent update -p cent-nuclei-templates -f
```
---
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
  - SOMETHING

# Files to exclude
exclude-files:
  - README.md
  - .gitignore
  - .pre-commit-config.yaml
  - LICENSE

# Add github urls
community-templates:
  - https://github.com/projectdiscovery/nuclei-templates
  - https://github.com/0x727/ObserverWard_0x727
  - https://github.com/0xAwali/Virtual-Host
  - https://github.com/0xPugazh/my-nuclei-templates
  - https://github.com/0xmaximus/final_freaking_nuclei_templates
  - https://github.com/1in9e/my-nuclei-templates
  - https://github.com/5cr1pt/templates
  - https://github.com/ARPSyndicate/kenzer-templates
  - https://github.com/Akokonunes/Private-Nuclei-Templates
  - https://github.com/Arvinthksrct/alltemplate
  - https://github.com/AshiqurEmon/nuclei_templates.git
  - https://github.com/CharanRayudu/Custom-Nuclei-Templates
  - https://github.com/DoubleTakes/nuclei-templates
  - https://github.com/Elsfa7-110/mynuclei-templates
  - https://github.com/ExpLangcn/NucleiTP
  - https://github.com/Harish4948/Nuclei-Templates
  - https://github.com/Hunt2behunter/nuclei-templates
  - https://github.com/Jagomeiister/nuclei-templates
  - https://github.com/JoshMorrison99/url-based-nuclei-templates
  - https://github.com/Kaue-Navarro/Templates-kaue-nuclei
  - https://github.com/KeepHowling/all_freaking_nuclei_templates
  - https://github.com/Lopseg/nuclei-c-templates
  - https://github.com/MR-pentestGuy/nuclei-templates
  - https://github.com/NightRang3r/misc_nuclei_templates
  - https://github.com/NitinYadav00/My-Nuclei-Templates
  - https://github.com/Odayex/Random-Nuclei-Templates
  - https://github.com/PedroFerreira97/nuclei_templates
  - https://github.com/R-s0n/Custom_Vuln_Scan_Templates
  - https://github.com/Saimonkabir/Nuclei-Templates
  - https://github.com/Saptak9983/Nuclei-Template
  - https://github.com/ShangRui-hash/my-nuclei-templates
  - https://github.com/freakyclown/Nuclei_templates
  - https://github.com/SirAppSec/nuclei-template-generator-log4j
  - https://github.com/Str1am/my-nuclei-templates
  - https://github.com/SumedhDawadi/Custom-Nuclei-Template
  - https://github.com/System00-Security/backflow
  - https://github.com/adampielak/nuclei-templates
  - https://github.com/aels/CVE-2022-37042
  - https://github.com/alexrydzak/rydzak-nuclei-templates
  - https://github.com/ayadim/Nuclei-bug-hunter
  - https://github.com/badboy-sft/badboy_17-Nuclei-Templates-Collection
  - https://github.com/binod235/nuclei-templates-and-reports
  - https://github.com/blazeinfosec/nuclei-templates
  - https://github.com/brinhosa/brinhosa-nuclei-templates
  - https://github.com/c-sh0/nuclei_templates
  - https://github.com/cipher387/juicyinfo-nuclei-templates
  - https://github.com/clarkvoss/Nuclei-Templates
  - https://github.com/coldrainh/nuclei-ByMyself
  - https://github.com/d3sca/Nuclei_Templates
  - https://github.com/daffainfo/my-nuclei-templates
  - https://github.com/damon-sec/Nuclei-templates-Collection
  - https://github.com/dk4trin/templates-nuclei
  - https://github.com/drfabiocastro/certwatcher-templates
  - https://github.com/ekinsb/Nuclei-Templates
  - https://github.com/erickfernandox/nuclei-templates
  - https://github.com/esetal/nuclei-bb-templates
  - https://github.com/ethicalhackingplayground/erebus-templates
  - https://github.com/geeknik/nuclei-templates-1
  - https://github.com/geeknik/the-nuclei-templates
  - https://github.com/glyptho/templatesallnuclei
  - https://github.com/im403/nuclei-temp
  - https://github.com/javaongsan/nuclei-templates
  - https://github.com/justmumu/SpringShell
  - https://github.com/kabilan1290/templates
  - https://github.com/kh4sh3i/CVE-2022-23131
  - https://github.com/luck-ying/Library-YAML-POC
  - https://github.com/mastersir-lab/nuclei-yaml-poc
  - https://github.com/mbskter/Masscan2Httpx2Nuclei-Xray
  - https://github.com/medbsq/ncl
  - https://github.com/meme-lord/Custom-Nuclei-Templates
  - https://github.com/n1f2c3/mytemplates
  - https://github.com/notnotnotveg/nuclei-custom-templates
  - https://github.com/obreinx/nuceli-templates
  - https://github.com/optiv/mobile-nuclei-templates
  - https://github.com/p0ch4t/nuclei-special-templates
  - https://github.com/panch0r3d/nuclei-templates
  - https://github.com/peanuth8r/Nuclei_Templates
  - https://github.com/pentest-dev/Profesional-Nuclei-Templates
  - https://github.com/pikpikcu/nuclei-templates
  - https://github.com/ping-0day/templates
  - https://github.com/praetorian-inc/chariot-launch-nuclei-templates
  - https://github.com/ptyspawnbinbash/template-enhancer
  - https://github.com/rafaelcaria/Nuclei-Templates
  - https://github.com/rafaelwdornelas/my-nuclei-templates
  - https://github.com/rahulkadavil/nuclei-templates
  - https://github.com/randomstr1ng/nuclei-sap-templates
  - https://github.com/redteambrasil/nuclei-templates
  - https://github.com/ree4pwn/my-nuclei-templates
  - https://github.com/ricardomaia/nuclei-template-generator-for-wordpress-plugins
  - https://github.com/sadnansakin/my-nuclei-templates
  - https://github.com/securitytest3r/nuclei_templates_work
  - https://github.com/sharathkramadas/k8s-nuclei-templates
  - https://github.com/shifa123/detections
  - https://github.com/sl4x0/NC-Templates
  - https://github.com/smaranchand/nuclei-templates
  - https://github.com/soapffz/myown-nuclei-poc
  - https://github.com/soumya123raj/Nuclei
  - https://github.com/souzomain/mytemplates
  - https://github.com/tamimhasan404/Open-Source-Nuclei-Templates-Downloader
  - https://github.com/test502git/log4j-fuzz-head-poc
  - https://github.com/th3r4id/nuclei-templates
  - https://github.com/thebrnwal/Content-Injection-Nuclei-Script
  - https://github.com/thecyberneh/nuclei-templatess
  - https://github.com/thelabda/nuclei-templates
  - https://github.com/themoonbaba/private_templates
  - https://github.com/rix4uni/BugBountyTips
  - https://github.com/Lu3ky13/Authorization-Nuclei-Templates
  - https://github.com/bug-vs-me/nuclei
  - https://github.com/topscoder/nuclei-wordfence-cve
  - https://github.com/toramanemre/apache-solr-log4j-CVE-2021-44228
  - https://github.com/toramanemre/log4j-rce-detect-waf-bypass
  - https://github.com/trickest/log4j
  - https://github.com/wasp76b/nuclei-templates
  - https://github.com/wr00t/templates
  - https://github.com/xinZa1/template
  - https://github.com/yarovit-developer/nuclei-templates
  - https://github.com/yavolo/nuclei-templates
  - https://github.com/z3bd/nuclei-templates
  - https://github.com/zer0yu/Open-PoC
  - https://github.com/zinminphyo0/KozinTemplates
  - https://github.com/bhataasim1/PersonalTemplates.git
  - https://github.com/pikpikcu/my-nuclei-templates
  - https://github.com/SirBugs/Priv8-Nuclei-Templates

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
