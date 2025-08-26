import { Card, Divider } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useParams } from "react-router";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import BasicInfoDescription from "@/views/cluster/detail/components/BasicInfoDescription.tsx";
import DetailExtra from "@/views/cluster/detail/components/DetailExtra.tsx";
import ConnectionDescription from "@/views/cluster/detail/components/ConnectionDescription.tsx";
import useApp from "antd/es/app/useApp";

export default function ClusterDetail() {
    const { t } = useLingui();
    const { id } = useParams();
    const { notification } = useApp();
    const { data, loading } = useRequest(
        () => {
            return clusterApi.clusterServiceGetCluster({ id: id as string });
        },
        {
            ready: !!id,
            refreshDeps: [id],
            onError: (e) => {
                notification.error({
                    message: t`get cluster error`,
                    description: e.message,
                });
            },
        },
    );
    return (
        <Card loading={loading} title={t`Cluster detail`} extra={<DetailExtra id={id} />}>
            <BasicInfoDescription cluster={data} />
            <Divider />
            <ConnectionDescription connection={data?.connection} />
        </Card>
    );
}
