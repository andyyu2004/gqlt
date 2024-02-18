# gqlt

### Tooling

#### Language Server

Install the language server `go install github.com/movio/gqlt/cmd/gqlt@latest`.

If your editor supports language servers, set it up appropriately.

#### VS Code integration

Build the language server as above.

```bash
git clone git@github.com:movio/gqlt
cd vscode-gqlt
npm i
npx vsce package
```

This outputs a file `gqlt-<version>.vsix`.
This can be installed from vscode using the command `Extensions: Install from VSIX`.

Set the following setting. It currently has to be the full path.
`"gqlt.path": "/home/<user>/go/bin/gqlt"`
