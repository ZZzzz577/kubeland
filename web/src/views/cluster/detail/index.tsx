import { Card, Divider } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useParams } from "react-router";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import BasicInfoDescription from "@/views/cluster/detail/components/description/BasicInfoDescription.tsx";
import DetailExtra from "@/views/cluster/detail/components/extra/DetailExtra.tsx";
import ConnectionDescription from "@/views/cluster/detail/components/description/ConnectionDescription.tsx";

export default function ClusterDetail() {
    const { t } = useLingui();
    const { id } = useParams();
    const { data, loading } = useRequest(
        () => {
            if (id) {
                return clusterApi.clusterServiceGetCluster({ id });
            }
            return Promise.resolve(undefined);
        },
        {
            ready: !!id,
            refreshDeps: [id],
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
