# sigsubs

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigsubs.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigsubs/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigsubs.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigsubs/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/License-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigsubs/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

sigsubs is a vertical correlation(subdomain discovery) tool. It gathers a list of subdomains passively using various online sources.

## Resources

* [Usage](#usage)
* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Post Installation](#post-installation)
* [Contribution](#contribution)

## Usage

To display help message for sigsubs use the `-h` flag:

```
$ sigsubs -h

     _                 _         
 ___(_) __ _ ___ _   _| |__  ___ 
/ __| |/ _` / __| | | | '_ \/ __|
\__ \ | (_| \__ \ |_| | |_) \__ \
|___/_|\__, |___/\__,_|_.__/|___/ V1.2.0
       |___/

USAGE:
  sigsubs [OPTIONS]

OPTIONS:
  -d                domain to find subdomains for
  -sE               comma separated list of sources to exclude
  -sL               list all the sources available
  -nC               no color mode
  -silent           silent mode: output subdomains only
  -sU               comma separated list of sources to use

```

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/drsigned/sigsubs/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

sigsubs requires go1.14+ to install successfully. Run the following command to get the repo

```bash
$ GO111MODULE=on go get -u -v github.com/drsigned/sigsubs/cmd/sigsubs
```

#### From Github

```bash
$ git clone https://github.com/drsigned/sigsubs.git; cd sigsubs/cmd/sigsubs/; go build; mv sigsubs /usr/local/bin/; sigsubs -h
```

## Post Installation

sigsubs will work after [installation](#installation). However, to configure sigsubs to work with certain services you will need to have setup API keys. Currently these services include:

* chaos
* github

The API keys are stored in the `$HOME/.config/sigsubs/conf.yaml` file - created upon first run - and uses the YAML format. Multiple API keys can be specified for each of these services.

Example:

```yaml
version: 1.2.0
sources:
    - alienvault
    - anubis
    - bufferover
    - cebaidu
    - certspotterv0
    - chaos
    - crtsh
    - github
    - hackertarget
    - rapiddns
    - riddler
    - sublist3r
    - threatcrowd
    - threatminer
    - urlscan
    - wayback
    - ximcx
keys:
    chaos:
        - d23a554bbc1aabb208c9acfbd2dd41ce7fc9db39asdsd54bbc1aabb208c9acfb
    github:
        - d23a554bbc1aabb208c9acfbd2dd41ce7fc9db39
        - asdsd54bbc1aabb208c9acfbd2dd41ce7fc9db39
```
## Contribution

[Issues](https://github.com/drsigned/sigsubs/issues) and [Pull Requests](https://github.com/drsigned/sigsubs/pulls) are welcome! 