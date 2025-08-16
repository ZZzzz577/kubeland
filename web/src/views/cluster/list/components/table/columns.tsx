import type { ColumnsType } from "antd/es/table";
import type { ApiV1ClusterCluster } from "@/generated";
import { Trans } from "@lingui/react/macro";
import { Space } from "antd";
import { Link } from "react-router";

export default function getClusterTableColumns(): ColumnsType<ApiV1ClusterCluster> {
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
            title: <Trans>description</Trans>,
            ellipsis: true,
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
