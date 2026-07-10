import "./wasm_exec.js";
import "@xterm/xterm/css/xterm.css";
import { FitAddon } from "@xterm/addon-fit";
import { Terminal } from "@xterm/xterm";
import lumingrate from "./lumingrate.wasm?url";

const terminalElement = document.getElementById("terminal");
const go = new Go();
const wasm = await WebAssembly.instantiateStreaming(fetch(lumingrate), go.importObject);
go.run(wasm.instance);

const term = new Terminal({
    fontSize: 14,
    lineHeight: 1.2,
    fontFamily: "Monaspace Neon, monospace",
    cursorBlink: true,
    theme: { background: "#000000", foreground: "#F7F7F2" },
});
const fit = new FitAddon();
term.loadAddon(fit);
term.open(terminalElement);
fit.fit();

globalThis.lumingrate.start({
    cols: term.cols,
    rows: term.rows,
    write: (data) => term.write(data),
});
term.onData((data) => globalThis.lumingrate.input(data));
term.onResize(({ cols, rows }) => globalThis.lumingrate.resize(cols, rows));

new ResizeObserver(() => {
    const dimensions = fit.proposeDimensions();
    if (dimensions) term.resize(dimensions.cols, dimensions.rows);
}).observe(terminalElement);

term.focus();
