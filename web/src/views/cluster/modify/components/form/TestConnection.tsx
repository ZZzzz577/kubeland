import { useLingui } from "@lingui/react/macro";
import { Button, Form, notification } from "antd";
import type { ApiV1ClusterCluster } from "@/generated";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";

export default function TestConnection() {
    const { t } = useLingui();
    const [notify, notifyContext] = notification.useNotification();
    const form = Form.useFormInstance<ApiV1ClusterCluster>();
    const { loading, run } = useRequest(clusterApi.clusterServiceTestConnection.bind(clusterApi), {
        manual: true,
        onSuccess: (data) => {
            notify.success({
                message: t`test connection success`,
                description: t`version is ${data.version}`,
            });
        },
        onError: (error) => {
            notify.error({
                message: t`test connection error`,
                description: error.message,
            });
            console.log(error);
        },
    });
    const handleTestConnection = () => {
        form.validateFields(["connection"], { recursive: true })
            .then((values) => {
                run({
                    apiV1ClusterConnection: values.connection,
                });
            })
            .catch(() => {});
    };
    return (
        <>
            {notifyContext}
            <Button type={"primary"} ghost loading={loading} onClick={handleTestConnection}>{t`test connection`}</Button>
        </>
    );
}