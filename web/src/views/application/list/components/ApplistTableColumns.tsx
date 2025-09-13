import type { ColumnsType } from "antd/es/table";
import type { ApiV1ApplicationApplication } from "@/generated";
import { Trans } from "@lingui/react/macro";
import { Link } from "react-router";
import { Space } from "antd";
import { ApplicationCluster } from "@/views/application/list/components/ApplicationListTable.tsx";
import EditButton from "@/views/application/list/components/EditButton.tsx";
import DeleteButton from "@/views/application/detail/components/DeleteButton.tsx";

export default function ApplicationTableColumns(): ColumnsType<ApiV1ApplicationApplication> {
    return [
        {
            dataIndex: "name",
            title: <Trans>Name</Trans>,
            render: (name: string) => {
                return <Link to={`/app/${name}`}>{name}</Link>;
            },
        },
        {
            title: <Trans>Cluster</Trans>,
            dataIndex: "clusterId",
            render: (clusterId?: string) => <ApplicationCluster id={clusterId} />,
        },
        {
            dataIndex: "description",
            title: <Trans>Description</Trans>,
            ellipsis: true,
        },
        {
            title: <Trans>Actions</Trans>,
            render: (_, record) => {
                return (
                    <Space>
                        <EditButton name={record.name} type={"link"} />
                        <DeleteButton name={record.name} type={"link"} />
                    </Space>
                );
            },
        },
    ];
}
