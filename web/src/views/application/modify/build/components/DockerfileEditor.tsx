import {type WrapperConfig} from 'monaco-editor-wrapper';
import {MonacoEditorReactComp} from '@typefox/monaco-editor-react';
import {configureDefaultWorkerFactory} from 'monaco-editor-wrapper/workers/workerLoaders';
import {LogLevel} from "@codingame/monaco-vscode-api/vscode/vs/platform/log/common/log";
import {useMemo} from "react";


export function DockerfileEditor() {
    const wrapperConfig = useMemo(() => createConfig(), []);

    return (
        <>
            <MonacoEditorReactComp
                wrapperConfig={wrapperConfig}
                style={{height: "100vh", width: "100vw"}}
            />
        </>
    )
}

const createConfig = (): WrapperConfig => {
    return {
        $type: 'extended',
        logLevel: LogLevel.Trace,
        extensions: [
            {
                config: {
                    name: "dockerfile",
                    publisher: "xxx",
                    version: '1.0.0',
                    engines: {
                        vscode: '*'
                    },
                    contributes: {
                        languages: [{
                            id: 'dockerfile',
                            extensions: ['.dockerfile', '.containerfile'],
                            filenames: ['Dockerfile', 'Containerfile'],
                            filenamePatterns: ["Dockerfile.*", "Containerfile.*"],
                            aliases: ["Docker", "Dockerfile", "Containerfile"]
                        }],
                    }
                }
            }
        ],
        editorAppConfig: {
            codeResources: {
                modified: {
                    uri: '/workspace/Dockerfile',
                    text: "FROM test",
                }
            },
            monacoWorkerFactory: configureDefaultWorkerFactory
        },
        languageClientConfigs: {
            configs: {
                dockerfile: {
                    name: "dockerfile",
                    connection: {
                        options: {
                            $type: "WebSocketUrl",
                            url: "ws://localhost:3000/ws"
                        }
                    },
                    clientOptions: {
                        documentSelector: ["dockerfile"],
                    },
                },
            }
        },
    };
}