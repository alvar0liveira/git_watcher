# Git Watcher
[![Go](https://github.com/alvar0liveira/git_watcher/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/alvar0liveira/git_watcher/actions/workflows/go.yml)


A command line utility for automatically commit changes to a repo.

Useful for Writing and keeping a log

To use create a env file where the repository is going to be:

```env
author_name="Your Name"
autor_email="Your Email"

```

Then use call command

```bash

git_watcher -cadence {0,1,2}

```

For help run ` git_watcher -help`
