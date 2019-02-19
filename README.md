[![Build Status](https://travis-ci.org/kairen/kubectl-config-merge.svg?branch=master)](https://travis-ci.org/kairen/kubectl-config-merge)

# kubectl-config-merge
This repository implements a kubectl plugin for merging multiple kubeconfig files.

![Screenshot](./screenshots/config-merge.gif)

## Usage
```
$ kubectl config-merge -h
Merge two or more kubeconfig files

Usage:
  config-merge [kubeconfig] [flags]

Examples:

	# Merge your kubeconfig and save it as "merge-config" in the current directory
	kubectl config-merge kubeconfig1 kubeconfig2 ...

	# Merge your kubeconfig and save it as "config" in "test" directory
	kubectl config-merge kubeconfig1 kubeconfig2 ... --path test/config

	# Merge your kubeconfig with $HOME/.kube/config as "$HOME/.kube/config"
	kubectl config-merge kubeconfig1 ... --home

	# To view merged kubeconfig result
	kubectl config-merge kubeconfig1 kubeconfig2 ... --view


Flags:
      --backup          If true, backup $HOME/.kube/config file to $HOME/.kube/config.bk (default true)
  -h, --help            help for config-merge
      --home            If true, merge with $HOME/.kube/config to $HOME/.kube/config
  -o, --output string   Merged kubeconfig output type(json or yaml) (default "yaml")
      --overwrite       If true, force merge kubeconfig
      --path string     Merged kubeconfig output name and path (default "merge-config")
      --version         Show config-merge version
      --view            View merged kubeconfig result
```