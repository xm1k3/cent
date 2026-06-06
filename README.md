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
</p>

# Install
```
go install -v github.com/xm1k3/cent/v2@latest
```

Or [download from releases](https://github.com/xm1k3/cent/releases)
<br><br>

after installation run `cent init` to initialize cent with the configuration files you find [here](https://github.com/xm1k3/cent/blob/main/.cent.yaml) 


# Supported commands

| Command | Description            |
| ------- | ---------------------- |
| check | Check if templates repo are still available |
| init    | Cent init configuration file      |
| summary | Print detailed summary of nuclei templates |
| update  | Update your repository |
| validate | Validate templates, if the template is invalid it is deleted from the folder |
| version  | Print cent version |

# Root flags

```
Flags:
      --config string   config file (default is .config/cent/.cent.yaml)
  -r, --by-repo         Group templates by repository (use with -k)
  -C, --console         Print console output
  -k, --keep-folders    Keep templates organized in folders by category
  -p, --path string     Root path to save the templates (default "cent-nuclei-templates")
  -t, --threads int     Number of threads to use when cloning repositories (default 10)
  -T, --timeout int     Timeout in minutes for each git clone (default 2)
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

## Basic Usage

Clone and insert all the community templates into the `cent-nuclei-templates` folder 
```
cent -p cent-nuclei-templates
```

Example output:
```
cent started
[CLONED] https://github.com/projectdiscovery/nuclei-templates
[CLONED] https://github.com/0xSojalSec/nuclei-templates-4
[CLONED] https://github.com/0xPugazh/my-nuclei-templates
[CLONED] https://github.com/0xSojalSec/my-nuclei-templates-1
[CLONED] https://github.com/0x727/ObserverWard
[CLONED] https://github.com/0xAwali/Blind-SSRF
[CLONED] https://github.com/0x727/ObserverWard_0x727
[CLONED] https://github.com/0xAwali/Virtual-Host
[CLONED] https://github.com/0xSojalSec/Nuclei-Templates-API-Linkfinder
...
... 
...
cent finished, you can find all your nuclei-templates in cent-nuclei-templates
```

### Keep templates organized by category
Use the `-k` flag to keep the internal folder structure of each repository:
```
cent -p cent-nuclei-templates -k
```
This will create a folder structure like:
```
cent-nuclei-templates/
  cves/
    CVE-2021-1234.yaml
  vulnerabilities/
    template.yaml
  others/
    uncategorized-template.yaml
```

Templates without a subfolder in their source repo are placed in `others/`.

### Group templates by repository
Use `-k -r` to also group templates by their source repository:
```
cent -p cent-nuclei-templates -k -r
```
```
cent-nuclei-templates/
  projectdiscovery/nuclei-templates/
    cves/
      CVE-2021-1234.yaml
    vulnerabilities/
      template.yaml
  user/repo-name/
    others/
      custom-template.yaml
```

## Summary Command

The `summary` command provides detailed statistics about your nuclei templates collection:

### Basic Summary
```bash
# Display summary in table format
cent summary

# Display summary in JSON format
cent summary --json
```

### Advanced Summary Features
```bash
# Limit number of tags displayed (default: 25)
cent summary --limit 10

# Search for specific data in summary
cent summary --search cve
cent summary --search wordpress
cent summary --search critical

# Update summary data
cent summary update

# Update with custom path
cent summary update -p /path/to/templates
```

### Summary Output Example
```
=== NUCLEI TEMPLATES SUMMARY ===

+-------------------+-------+
| METRIC            | COUNT |
+-------------------+-------+
| Total Templates   |  3249 |
| CVE Templates     |  3821 |
| Invalid Templates |     1 |
| Valid Templates   |  3248 |
+-------------------+-------+

=== SEVERITY DISTRIBUTION ===
+----------+-------+
| SEVERITY | COUNT |
+----------+-------+
| CRITICAL |   582 |
| HIGH     |   877 |
| MEDIUM   |   877 |
| LOW      |    63 |
| INFO     |   744 |
+----------+-------+

=== TOP TAGS ===
+---------------+-------+
| TAG           | COUNT |
+---------------+-------+
| cve           |  1909 |
| xss           |   569 |
| wordpress     |   487 |
| lfi           |   459 |
| wp-plugin     |   450 |
+---------------+-------+
```

### JSON Output Structure
```json
{
  "metrics": {
    "total_templates": 3249,
    "cve_templates": 3821,
    "invalid_templates": 1,
    "valid_templates": 3248
  },
  "severity_distribution": {
    "CRITICAL": 582,
    "HIGH": 877,
    "MEDIUM": 877,
    "LOW": 63,
    "INFO": 744
  },
  "tags": {
    "cve": 1909,
    "xss": 569,
    "wordpress": 487
  },
  "last_updated": "2024-01-15 14:30:25"
}
```

## Update Command

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

Example output:
```
[D][-] Dir  removed	cent-nuclei-templates/dns
[D][-] Dir  removed	cent-nuclei-templates/dns/subdomain
```

Same thing with `exclude-files`
```
cent update -p cent-nuclei-templates -f
```

## Configuration Management

### Initialize Configuration
```bash
# Initialize with default configuration
cent init

# Initialize with custom URL
cent init --url https://example.com/config.yaml

# Overwrite existing configuration
cent init --overwrite
```

### Check Configuration Status
```bash
# Check if configuration file exists
cent init check
```

### Check Template Repositories
```bash
# Check if all template repositories are accessible
cent check

# Remove inaccessible repositories from config
cent check --remove
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
You need to configure `cent` parameters in `.config/cent/.cent.yaml`
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

# Add github urls (simple format)
community-templates:
  - https://github.com/projectdiscovery/nuclei-templates
  - https://github.com/other/repo

  # Extended format: pin a repo to a specific commit
  - url: https://github.com/user/repo
    commit: abc123f

  # Extended format: per-repo exclude
  - url: https://github.com/user/repo2
    exclude:
      - "workflows/"
      - "fuzzing/"
```

Both formats can be mixed in the same config file. When `commit` is specified, cent will clone the full repo and checkout that specific commit instead of using a shallow clone. Per-repo `exclude` filters out matching paths from that specific repository.

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
