import {useLingui} from "@lingui/react/macro";
import {Button, notification} from "antd";
import type {ApiV1ClusterConnection} from "@/generated";
import {useRequest} from "ahooks";
import {clusterApi} from "@/api";

export default function TestConnection(props: { connection?: ApiV1ClusterConnection }) {
    const {t} = useLingui();
    const {connection} = props;
    const [notify, notifyContext] = notification.useNotification();
    const {loading, run} = useRequest(clusterApi.clusterServiceTestConnection.bind(clusterApi), {
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
        if (connection) {
            run({
                apiV1ClusterConnection: connection,
            });
        }
    };
    return (
        <>
            {notifyContext}
            <Button type={"primary"} ghost loading={loading}
                    onClick={handleTestConnection}>{t`Test connection`}</Button>
        </>
    );
}