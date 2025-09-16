import { Button, Card, Popconfirm, Space, Table } from "antd";
import { usePagination, useRequest } from "ahooks";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { imageApi } from "@/api";
import type { ApiV1ImageImageRepo } from "@/generated";
import { PlusOutlined } from "@ant-design/icons";
import { Link, useNavigate } from "react-router";
import type { ColumnsType } from "antd/es/table";

export default function ImageRepoList() {
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { loading, pagination, data, refresh } = usePagination(
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

    const { run: deleteImage, loading: deleteLoading } = useRequest(
        imageApi.imageServiceDeleteImageRepo.bind(imageApi),
        {
            manual: true,
            onSuccess: () => {
                notification.success({
                    message: t`delete image repository success`,
                });
                refresh();
            },
            onError: (e) => {
                notification.error({
                    message: t`failed to delete image repository`,
                    description: e.message,
                });
            },
        },
    );
    const handleDelete = (name?: string) => {
        if (name) {
            deleteImage({ name });
        }
    };

    const columns: ColumnsType<ApiV1ImageImageRepo> = [
        {
            title: t`Name`,
            dataIndex: "name",
            render: (name: string) => {
                return <Link to={`/image/${name}`}>{name}</Link>;
            },
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
        {
            title: t`Action`,
            render: (_, record) => {
                return (
                    <Space>
                        <Button type={"link"} onClick={() => navigate(`/image/${record.name}/edit`)}>
                            Edit
                        </Button>
                        <Popconfirm
                            title={t`Are you sure to delete this image repository?`}
                            onConfirm={() => handleDelete(record.name)}
                        >
                            <Button type={"link"} danger loading={deleteLoading}>
                                Delete
                            </Button>
                        </Popconfirm>
                    </Space>
                );
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