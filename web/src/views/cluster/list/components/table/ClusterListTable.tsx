import { usePagination } from "ahooks";
import { clusterApi } from "@/api";
import { Table } from "antd";
import getClusterTableColumns from "@/views/cluster/list/components/table/columns.tsx";

export default function ClusterListTable() {
    const { loading, pagination, data } = usePagination(async ({ current, pageSize }) => {
        const res = await clusterApi.clusterServiceListClusters({ pageCurrent: current, pageSize });
        return {
            current: res.pagination?.current ?? 0,
            total: res.pagination?.total ?? 0,
            list: res.items ?? [],
        };
    });
    const columns = getClusterTableColumns();
    return <Table rowKey={"id"} loading={loading} pagination={pagination} columns={columns} dataSource={data?.list} />;
}
