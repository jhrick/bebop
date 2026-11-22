# Bop

Bop is a CLI tool that runs [OpenCode](https://opencode.ai) inside Docker containers, providing an isolated and reusable development environment for AI-assisted coding.

It keeps the execution environment and its configuration separate from individual project repositories, allowing the same setup to be reused across different projects.

## Goals

- Run OpenCode in isolated Docker containers.
- Keep development environment configuration outside project repositories.
- Reuse the same environment across multiple projects.
- Simplify container lifecycle management through a CLI.
