how-ya-doing
============

[![License](https://img.shields.io/badge/licence-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0) [![Build Status](https://travis-ci.org/blacktop/how-ya-doing.svg?branch=master)](https://travis-ci.org/blacktop/how-ya-doing)

Monitor your Github repos activity

Getting Started
---------------

```bash
$ brew install blacktop/tap/hyd
```

Next export a `GITHUB_ACCESS_TOKEN` environment variable with the `repo` scope selected. This will be used to OAuth to the Github API. Create yours [here](https://github.com/settings/tokens/new).

```console
$ export GITHUB_ACCESS_TOKEN=`YOUR_TOKEN`
```

```console
$ hyd maliceio/malice
```

![screen](https://github.com/blacktop/how-ya-doing/raw/master/screen-shot.png)

> To quit press `Ctrl+c` or just `q`

TODO
----

-	[ ] add database for history beyond what Github gives (2 weeks)
