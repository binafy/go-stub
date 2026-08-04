# Security Policy

## Supported versions

go-stub is pre-1.0; security fixes land on the latest tagged release and `main`.

## Reporting a vulnerability

Please **do not** open a public issue for security problems. Instead, use
GitHub's [private vulnerability reporting](https://github.com/binafy/go-stub/security/advisories/new)
to disclose the issue privately.

Include a description, a reproducer if possible, and the affected version. We aim
to acknowledge reports within a few days.

## Scope notes

go-stub reads template files and writes generated files on the local filesystem.
When you feed it untrusted stub paths or destination paths, treat path handling
as you would any file I/O: validate inputs, since a malicious path or placeholder
value can direct output outside an intended directory.
