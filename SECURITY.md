# Security Policy

## Supported versions

The latest tagged release receives security fixes. arsenal is pre-1.0; pin a
specific version for reproducible engagements.

## Reporting a vulnerability

Please report security issues privately rather than opening a public issue. Use
GitHub's private vulnerability reporting ("Report a vulnerability" under the
Security tab) for this repository. Include:

- a description of the issue and its impact,
- steps to reproduce or a proof of concept,
- affected version (`arsenal version`).

You can expect an acknowledgement within a few days.

## Scope and intended use

arsenal installs and isolates offensive-security tooling for authorized testing,
research, and education. The registry pins upstream tools but does not audit
their behavior; review the tools you install. Running tools may require elevated
privileges and can affect networks and hosts - only use them where you have
explicit authorization.

## Hardening notes

- arsenal makes network calls only for registry sync and tool installs.
- It does not collect telemetry.
- Installs are isolated per tool/version; removing a tool does not affect others.
