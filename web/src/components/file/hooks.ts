import { useState } from "react";
import { useLingui } from "@lingui/react/macro";

export function useFileUpload(options?: {
    onSuccess?: (file: File, content: string) => void;
    onError?: (error: Error) => void;
}) {
    const { t } = useLingui();
    const [loading, setLoading] = useState<boolean>(false);

    const readFileContent = (file: File): Promise<string> => {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = (event) => {
                if (typeof event.target?.result === "string") {
                    resolve(event.target.result);
                } else {
                    reject(new Error(t`unable to read the file content`));
                }
            };
            reader.onerror = () => reject(reader.error);
            reader.readAsText(file);
        });
    };

    const handleFileUpload = async (file: File) => {
        setLoading(true);
        try {
            const content = await readFileContent(file);
            options?.onSuccess?.(file, content);
        } catch (error: unknown) {
            if (error instanceof Error) {
                options?.onError?.(error);
            } else {
                // 处理非 Error 类型的错误
                options?.onError?.(new Error(String(error)));
            }
        } finally {
            setLoading(false);
        }
        return false;
    };
    return {
        loading,
        handleFileUpload,
    };
}
