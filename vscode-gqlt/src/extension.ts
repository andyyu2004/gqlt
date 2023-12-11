import * as vscode from "vscode";
import { workspace } from "vscode";
import * as lc from "vscode-languageclient/node";
import { bootstrap } from "./bootstrap";
import { makeConfig } from "./config";

interface LspContext {
  client: lc.LanguageClient;
  subscriptions: vscode.Disposable[];
}
let lcx: LspContext | undefined;

export async function activate(context: vscode.ExtensionContext) {
  context.subscriptions.push(
    vscode.commands.registerCommand("gqlt.restart-server", async () => {
      await vscode.window.showInformationMessage("Restarting gqlt...");
      deactivate();
      while (context.subscriptions.length > 0) {
        context.subscriptions.pop()!.dispose();
      }

      await activateInner(context);
    })
  );

  await activateInner(context);
}

async function activateInner(context: vscode.ExtensionContext) {
  const config = makeConfig();
  const serverPath = await bootstrap(context, config);
  console.log("running gqlt server at", serverPath);

  const opt: lc.Executable = {
    command: serverPath,
    options: { env: config.env },
  };
  const serverOptions: lc.ServerOptions = {
    run: opt,
    debug: opt,
  };

  const clientOptions: lc.LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "gqlt" }],
  };

  const client = new lc.LanguageClient(
    "gqlt",
    "gqlt",
    serverOptions,
    clientOptions
  );

  lcx = { client, subscriptions: context.subscriptions };
  context.subscriptions.push(client.start());
}

export function deactivate() {
  lcx?.client.stop();
  lcx = undefined;
}
