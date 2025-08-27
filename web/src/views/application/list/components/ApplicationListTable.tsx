import { usePagination, useRequest } from "ahooks";
import { applicationApi, clusterApi } from "@/api";
import { Spin, Table } from "antd";
import ApplicationTableColumns from "@/views/application/list/components/ApplistTableColumns.tsx";
import useApp from "antd/es/app/useApp";
import { useLingui } from "@lingui/react/macro";

export default function ApplicationListTable() {
    const { t } = useLingui();
    const { notification } = useApp();
    const { loading, pagination, data } = usePagination(
        async ({ current, pageSize }) => {
            const res = await applicationApi.applicationServiceListApplications({ pageCurrent: current, pageSize });
            return {
                current: res.pagination?.current ?? 0,
                total: res.pagination?.total ?? 0,
                list: res.items ?? [],
            };
        },
        {
            onError: (e) => {
                notification.error({
                    message: t`failed to get application list`,
                    description: e.message,
                });
            },
        },
    );
    const columns = ApplicationTableColumns();
    return <Table rowKey={"id"} loading={loading} pagination={pagination} columns={columns} dataSource={data?.list} />;
}

export function ApplicationCluster(props: { id?: string }) {
    const { data: clusterList, loading: clusterLoading } = useRequest(
        clusterApi.clusterServiceListClusters.bind(clusterApi),
    );
    return (
        <Spin spinning={clusterLoading}>{clusterList?.items?.find((item) => item.id === props.id)?.name || ""}</Spin>
    );
}