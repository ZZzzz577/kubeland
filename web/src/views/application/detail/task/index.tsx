import { useNavigate, useParams } from "react-router";
import { useRequest } from "ahooks";
import { buildTaskApi } from "@/api";
import useApp from "antd/es/app/useApp";
import { useLingui } from "@lingui/react/macro";
import { Button, Card, Space, Table, Tag } from "antd";
import { type ApiV1BuildTaskBuildTask, ApiV1BuildTaskBuildTaskStatusEnum } from "@/generated";
import type { ColumnsType } from "antd/es/table";
import { CheckCircleOutlined, CloseCircleOutlined, ExclamationCircleOutlined, SyncOutlined } from "@ant-design/icons";
import CreateBuildTask from "@/views/application/detail/task/components/CreateBuildTask.tsx";

export default function BuildTask() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { data, loading, refresh } = useRequest(buildTaskApi.buildTaskServiceList.bind(buildTaskApi), {
        ready: !!name,
        defaultParams: [{ name: name as string }],
        onError: (e) =>
            notification.error({
                message: t`failed to get build task list`,
                description: e.message,
            }),
    });

    const columns: ColumnsType<ApiV1BuildTaskBuildTask> = [
        {
            title: t`name`,
            dataIndex: "name",
        },
        {
            title: t`status`,
            dataIndex: "status",
            render: (value?: ApiV1BuildTaskBuildTaskStatusEnum) => {
                switch (value) {
                    case ApiV1BuildTaskBuildTaskStatusEnum.StatusRunning:
                        return (
                            <Tag icon={<SyncOutlined spin />} color="processing">
                                Running
                            </Tag>
                        );
                    case ApiV1BuildTaskBuildTaskStatusEnum.StatusComplete:
                    case ApiV1BuildTaskBuildTaskStatusEnum.StatusSuccessCriteriaMet:
                        return (
                            <Tag icon={<CheckCircleOutlined />} color="success">
                                Complete
                            </Tag>
                        );
                    case ApiV1BuildTaskBuildTaskStatusEnum.StatusFailed:
                    case ApiV1BuildTaskBuildTaskStatusEnum.StatusFailureTarget:
                        return (
                            <Tag icon={<CloseCircleOutlined />} color="error">
                                Failed
                            </Tag>
                        );
                    case ApiV1BuildTaskBuildTaskStatusEnum.StatusSuspended:
                        return (
                            <Tag icon={<ExclamationCircleOutlined />} color="warning">
                                Suspended
                            </Tag>
                        );
                    default:
                        return <Tag color="default">Unknown</Tag>;
                }
            },
        },
        {
            title: t`created at`,
            dataIndex: "createdAt",
            render: (value: Date) => value.toLocaleString(),
        },
        {
            title: "",
            dataIndex: "name",
            render: (value: string) => (
                <Space>
                    <Button type={"link"} onClick={() => navigate(`/app/${name}/task/${value}`)}>{t`Detail`}</Button>
                </Space>
            ),
        },
    ];
    return (
        <Card
            title={t`Build tasks`}
            className={"!rounded-t-none"}
            variant="borderless"
            loading={loading}
            extra={
                <Space>
                    <CreateBuildTask />
                    <Button icon={<SyncOutlined />} loading={loading} onClick={refresh}>{t`Refresh`}</Button>
                </Space>
            }
        >
            <Table<ApiV1BuildTaskBuildTask>
                rowKey={"name"}
                loading={loading}
                columns={columns}
                dataSource={data?.items}
            />
        </Card>
    );
}