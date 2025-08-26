import { useLingui } from "@lingui/react/macro";
import { Button, Form } from "antd";
import type { ApiV1ClusterCluster } from "@/generated";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import useApp from "antd/es/app/useApp";

export default function TestConnection() {
    const { t } = useLingui();
    const { notification } = useApp();
    const form = Form.useFormInstance<ApiV1ClusterCluster>();
    const { loading, run } = useRequest(clusterApi.clusterServiceTestConnection.bind(clusterApi), {
        manual: true,
        onSuccess: (data) => {
            notification.success({
                message: t`test connection success`,
                description: t`version is ${data.version}`,
            });
        },
        onError: (error) => {
            notification.error({
                message: t`test connection error`,
                description: error.message,
            });
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
            <Button type={"primary"} ghost loading={loading} onClick={handleTestConnection}>{t`test connection`}</Button>
        </>
    );
}