import { useEffect, useRef } from "react";
import { Terminal } from "@xterm/xterm";
import { AttachAddon } from "@xterm/addon-attach";
import { FitAddon } from "@xterm/addon-fit";
import { WebLinksAddon } from "@xterm/addon-web-links";
import "@xterm/xterm/css/xterm.css";

export default function LogView(props: { className?: string; url: string }) {
    const { className, url } = props;
    const terminalRef = useRef(null);
    useEffect(() => {
        if (!terminalRef.current) return;
        const terminal = new Terminal({
            convertEol: true,
            fontSize: 16,
            fontWeight: "500",
            theme: { background: "black" },
        });

        const fitAddon = new FitAddon();
        terminal.loadAddon(fitAddon);

        terminal.loadAddon(new WebLinksAddon());

        const webSocket = new WebSocket(url);
        const attachAddon = new AttachAddon(webSocket);
        terminal.loadAddon(attachAddon);

        terminal.open(terminalRef.current);

        const resizeObserver = new ResizeObserver(() => {
            fitAddon.fit();
        });
        resizeObserver.observe(terminalRef.current);

        return () => {
            resizeObserver.disconnect();
            webSocket.close();
            terminal.dispose();
        };
    }, [url]);
    return <div className={className} ref={terminalRef} />;
}
