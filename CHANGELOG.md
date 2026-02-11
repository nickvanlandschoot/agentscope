# Changelog

All notable changes to AgentScope will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-02-10

### Added

- Portal directory system with symlink support for isolated Claude sessions
- Automatic cleanup of temporary portal directories on session exit
- Support for passing Claude CLI arguments through `activate` command

## [0.1.0] - 2024-02-08

### Added

- `init` command to initialize `.agentscope` directory
- `new` command to create instruction files
- `activate` command with interactive selection
- YAML frontmatter support for instruction files
- Position-based sorting of instructions
- Default enablement configuration
- Dynamic CLAUDE.md generation
