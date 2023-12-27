import { fileURLToPath } from 'node:url';
import { fork } from 'node:child_process';

const cliEntry = new URL("../dist/cli-run.mjs", import.meta.url);
if (process.argv[2] === "dev") {
  process.env.__CLI_ARGV__ = JSON.stringify(process.argv);
  startSubprocess();
} else {
  import(cliEntry.href);
}
function startSubprocess() {
  let childProc;
  const onShutdown = () => {
    if (childProc) {
      childProc.kill();
      childProc = void 0;
    }
  };
  process.on("exit", onShutdown);
  process.on("SIGTERM", onShutdown);
  process.on("SIGINT", onShutdown);
  process.on("SIGQUIT", onShutdown);
  start();
  function start() {
    const _argv = (process.env.__CLI_ARGV__ ? JSON.parse(process.env.__CLI_ARGV__) : process.argv).slice(2);
    const execArguments = getInspectArgs();
    childProc = fork(fileURLToPath(cliEntry), [], { execArgv: execArguments });
    childProc.on("close", (code) => {
      if (code) {
        process.exit(code);
      }
    });
    childProc.on("message", (message) => {
      if (message?.type === "nuxt:restart") {
        childProc?.kill();
        startSubprocess();
      }
    });
    function getInspectArgs() {
      const inspectArgv = _argv.find((argvItem) => argvItem.includes("--inspect"));
      if (!inspectArgv) {
        return [];
      }
      return [inspectArgv];
    }
  }
}
