# lazyworkflows

CLI tool for managing GitHub workflows. Heavily inspired by lazygit.

This tool is at its very early stages. This project is a way for me to try out different things with Rust, but hopefully become a tool that would be suited actual use.



## TODOs


Things that I'd like to do, loosely formulated:

- **WORKING ON THIS:** Have some sort of DSL that can communicate the responses from the GitHub API in a way that we can work with in the program
- Have the app download each workflow file and figure out the inputs to the workflow file
  - Will have to be able to understand YAML
  - Has to look at `workflow_dispatch` trigger to determine them
  - Somehow let users be able to configure the inputs before passing them to the dispatch function
- Create a GUI
  - Should be TUI based like lazygit.
  - Should be able to navigate between workflows
  - Layout should respect workflows based on owners/repos
  - Layout should be able to give an overview of running workflows
  - Layout should give an indication if a workflow is running, pending(?) and finished (with added visual effects for either successes or failures)

