import type {ColumnsType} from "antd/es/table";
import type {ApiV1ApplicationApplication} from "@/generated";
import {Trans} from "@lingui/react/macro";
import {Link} from "react-router";
import {Space} from "antd";
import ApplicationDeleteButton from "@/views/application/list/components/ApplicationDeleteButton.tsx";
import {ApplicationCluster} from "@/views/application/list/components/ApplicationListTable.tsx";

export default function ApplicationTableColumns(): ColumnsType<ApiV1ApplicationApplication> {
    return [
        {
            dataIndex: "name",
            title: <Trans>Name</Trans>,
            render: (name: string, record) => {
                return <Link to={`/app/application/${record.id}/detail`}>{name}</Link>;
            },
        },
        {
            title: <Trans>Cluster</Trans>,
            dataIndex: "clusterId",
            render: (clusterId?: string) => <ApplicationCluster id={clusterId}/>,
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
                        <Link to={`/app/application/${record.id}/edit`}>
                            <Trans>edit</Trans>
                        </Link>
                        <ApplicationDeleteButton id={record.id}/>
                    </Space>
                );
            },
        },
    ];
}