import type { ColumnsType } from "antd/es/table";
import type { ApiV1ClusterCluster } from "@/generated";
import { Trans } from "@lingui/react/macro";
import { Space } from "antd";
import { Link } from "react-router";
import { ClusterStatus, OperatorStatus } from "@/views/cluster/list/components/ClusterListTable.tsx";

export default function ClusterListTableColumns(): ColumnsType<ApiV1ClusterCluster> {
    return [
        {
            dataIndex: "name",
            title: <Trans>Name</Trans>,
            render: (name: string, record) => {
                return <Link to={`/cluster/${record.id}/detail`}>{name}</Link>;
            },
        },
        {
            dataIndex: "description",
            title: <Trans>Description</Trans>,
            ellipsis: true,
        },
        {
            title: <Trans>Address</Trans>,
            dataIndex: ["connection", "address"],
            ellipsis: true,
        },
        {
            title: <Trans>Connection status</Trans>,
            render: (_: string, record) => <ClusterStatus conn={record.connection} />,
        },
        {
            title: <Trans>Operator status</Trans>,
            render: (_: string, record) => <OperatorStatus conn={record.connection} />,
        },
        {
            title: <Trans>Create time</Trans>,
            dataIndex: "createdAt",
            render: (createdAt: Date) => {
                return createdAt.toLocaleString();
            },
        },
        {
            title: <Trans>Update time</Trans>,
            dataIndex: "updatedAt",
            render: (updatedAt: Date) => {
                return updatedAt.toLocaleString();
            },
        },
        {
            title: <Trans>Actions</Trans>,
            render: (_, record) => {
                return (
                    <Space>
                        <Link to={`/cluster/${record.id}/edit`}>
                            <Trans>edit</Trans>
                        </Link>
                    </Space>
                );
            },
        },
    ];
}
