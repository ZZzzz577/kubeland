import { Button, Card, Space, Table } from "antd";
import { usePagination } from "ahooks";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { imageApi } from "@/api";
import type { ApiV1ImageImageRepo } from "@/generated";
import { PlusOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router";
import type { ColumnsType } from "antd/es/table";

export default function ImageRepoList() {
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { loading, pagination, data } = usePagination(
        async ({ current, pageSize }) => {
            const res = await imageApi.imageServiceListImageRepos({ pageCurrent: current, pageSize });
            return {
                current: res.pagination?.current ?? 0,
                total: res.pagination?.total ?? 0,
                list: res.items ?? [],
            };
        },
        {
            onError: (e) => {
                notification.error({
                    message: t`failed to get image repository list`,
                    description: e.message,
                });
            },
        },
    );
    const columns: ColumnsType<ApiV1ImageImageRepo> = [
        {
            title: t`Name`,
            dataIndex: "name",
        },
        {
            title: t`Description`,
            dataIndex: "description",
            ellipsis: true,
        },
        {
            title: t`URL`,
            dataIndex: "url",
            ellipsis: true,
        },
        {
            title: t`Create time`,
            dataIndex: "createdAt",
            render: (text) => {
                return text.toLocaleString();
            },
        },
        {
            title: t`Update time`,
            dataIndex: "updatedAt",
            render: (text) => {
                return text.toLocaleString();
            },
        },
    ];
    return (
        <Card
            extra={
                <Space>
                    <Button icon={<PlusOutlined />} type={"primary"} onClick={() => navigate("/image/create")}>
                        {t`Create`}
                    </Button>
                </Space>
            }
        >
            <Table<ApiV1ImageImageRepo>
                rowKey={"name"}
                loading={loading}
                columns={columns}
                dataSource={data?.list}
                pagination={pagination}
            />
        </Card>
    );
}
