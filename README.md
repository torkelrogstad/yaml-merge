`yaml-merge` lets you merge YAML files. Neat!

The merging is rather crude:

1. Everything except maps and lists (lists are concatenated) are just crudely
   overwritten, by whatever comes last in the list of YAML files.
2. Any formatting you had is completely wiped out post-processing.
3. There are no flags to control any of this behavior (yet).

## Build/installation instructions:

> If have to ask, you're not ready.

This is quickly-whipped out code with no tests. There's probably bugs and other
issues here. Use at your own risk!
