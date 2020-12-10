# sigsubs

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigsubs.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigsubs/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigsubs.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigsubs/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/License-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigsubs/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

sigsubs is a vertical correlation(subdomain discovery) tool. It gathers a list of subdomains passively using various online sources.

## Resources

* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Usage](#usage)
* [Contribution](#contribution)

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

## Usage

To display help message for sigsubs use the `-h` flag:

```
$ sigsubs -h

     _                 _         
 ___(_) __ _ ___ _   _| |__  ___ 
/ __| |/ _` / __| | | | '_ \/ __|
\__ \ | (_| \__ \ |_| | |_) \__ \
|___/_|\__, |___/\__,_|_.__/|___/ V1.0.0
       |___/

USAGE:
  sigsubs [OPTIONS]

OPTIONS:
  -d                domain to find subdomains for
  -exclude          comma separated list of sources to exclude
  -ls               list all the sources available
  -nc               no color mode
  -s                silent mode: output subdomains only
  -use              comma separated list of sources to use
```

## Contribution

[Issues](https://github.com/drsigned/sigsubs/issues) and [Pull Requests](https://github.com/drsigned/sigsubs/pulls) are welcome! 