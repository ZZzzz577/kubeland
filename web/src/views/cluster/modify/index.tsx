import { Card } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useParams } from "react-router";
import ClusterModifyForm from "@/views/cluster/modify/components/ClusterModifyForm.tsx";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import BackButton from "@/components/back/BackButton.tsx";
import useApp from "antd/es/app/useApp";

export default function ClusterModify() {
    const { t } = useLingui();
    const { notification } = useApp();
    const { id } = useParams();
    const isUpdate = !!id;
    const title = isUpdate ? t`Update cluster` : t`Create cluster`;

    const { data, loading } = useRequest(
        () => {
            if (id) {
                return clusterApi.clusterServiceGetCluster({ id });
            }
            return Promise.resolve(undefined);
        },
        {
            ready: isUpdate,
            refreshDeps: [id],
            onError(e) {
                notification.error({
                    message: t`failed to get cluster`,
                    description: e.message,
                });
            },
        },
    );

    return (
        <Card title={title} loading={loading} extra={<BackButton />}>
            <ClusterModifyForm cluster={data} />
        </Card>
    );
}
