SHELL := /usr/bin/env bash

GOLIB ?= golib

.PHONY: check check-full ci cohesion config inventory repository-check specification-check specification-check-online workflows

config:
	$(GOLIB) config validate

check:
	$(GOLIB) check --local --all

check-full: specification-check-online
	$(GOLIB) check --all

ci: config inventory repository-check cohesion specification-check workflows check

cohesion:
	$(GOLIB) cohesion check

inventory:
	$(GOLIB) inventory

repository-check:
	$(GOLIB) repository check

specification-check:
	$(GOLIB) specification check

specification-check-online:
	$(GOLIB) specification check --online

workflows:
	$(GOLIB) workflows check
