import { ExtensionContext, workspace } from "vscode";

export interface Config {
  path?: string;
  env?: object;
}

export function makeConfig(): Config {
  const cfg = workspace.getConfiguration("gqlt");
  return {
    path: cfg.get("path"),
    env: cfg.get("env"),
  };
}
