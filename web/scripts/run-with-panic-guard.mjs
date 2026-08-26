import { spawn } from "node:child_process";

const invocation = process.argv.slice(2);
if (invocation[0] === "--") invocation.shift();
const [command, ...args] = invocation;

if (!command) {
    console.error("panic guard requires a command after --");
    process.exit(2);
}

const panicPattern = /Rolldown panicked|thread ['\"]?rolldown[^\r\n]* panicked|ModuleLoader[^\r\n]*main thread/i;
let captured = "";

const child = spawn(command, args, {
    env: process.env,
    shell: false,
    stdio: ["inherit", "pipe", "pipe"],
    windowsHide: true,
});

function forward(stream, destination) {
    stream.on("data", (chunk) => {
        destination.write(chunk);
        captured = `${captured}${chunk.toString()}`.slice(-512_000);
    });
}

forward(child.stdout, process.stdout);
forward(child.stderr, process.stderr);

const { code, signal } = await new Promise((resolve, reject) => {
    child.once("error", reject);
    child.once("close", (exitCode, exitSignal) => resolve({ code: exitCode, signal: exitSignal }));
});

if (panicPattern.test(captured)) {
    console.error("panic guard rejected a Rolldown worker panic even though the child process may have exited successfully");
    process.exit(86);
}
if (signal) {
    console.error(`guarded command terminated by signal ${signal}`);
    process.exit(1);
}
process.exit(code ?? 1);
