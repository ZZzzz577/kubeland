import { useLingui } from "@lingui/react/macro";
import { Button } from "antd";
import type { ApiV1ClusterConnection } from "@/generated";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import useApp from "antd/es/app/useApp";

export default function TestConnection(props: { connection?: ApiV1ClusterConnection }) {
    const { t } = useLingui();
    const { connection } = props;
    const { notification } = useApp();
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
        if (connection) {
            run({
                apiV1ClusterConnection: connection,
            });
        }
    };
    return (
        <>
            <Button
                type={"primary"}
                ghost
                loading={loading}
                onClick={handleTestConnection}
            >{t`Test connection`}</Button>
        </>
    );
}
