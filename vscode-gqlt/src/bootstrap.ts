import * as vscode from "vscode";
import { Config } from "./config";
import assert = require("assert");

export async function bootstrap(
  context: vscode.ExtensionContext,
  config: Config
): Promise<string> {
  return serverPath(context, config);
}

async function serverPath(
  context: vscode.ExtensionContext,
  config: Config
): Promise<string> {
  if (config.path) {
    return config.path;
  }

  return "gqlt";
}
