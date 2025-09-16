import { Button, Card, Popconfirm, Space, Table } from "antd";
import { usePagination, useRequest } from "ahooks";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { gitApi } from "@/api";
import type { ApiV1GitGitRepo } from "@/generated";
import { PlusOutlined } from "@ant-design/icons";
import { Link, useNavigate } from "react-router";
import type { ColumnsType } from "antd/es/table";

export default function GitRepoList() {
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { loading, pagination, data, refresh } = usePagination(
        async ({ current, pageSize }) => {
            const res = await gitApi.gitServiceListGitRepos({ pageCurrent: current, pageSize });
            return {
                current: res.pagination?.current ?? 0,
                total: res.pagination?.total ?? 0,
                list: res.items ?? [],
            };
        },
        {
            onError: (e) => {
                notification.error({
                    message: t`failed to get git repository list`,
                    description: e.message,
                });
            },
        },
    );

    const { run: deleteGit, loading: deleteLoading } = useRequest(gitApi.gitServiceDeleteGitRepo.bind(gitApi), {
        manual: true,
        onSuccess: () => {
            notification.success({
                message: t`delete git repository success`,
            });
            refresh();
        },
        onError: (e) => {
            notification.error({
                message: t`failed to delete git repository`,
                description: e.message,
            });
        },
    });
    const handleDelete = (name?: string) => {
        if (name) {
            deleteGit({ name });
        }
    };

    const columns: ColumnsType<ApiV1GitGitRepo> = [
        {
            title: t`Name`,
            dataIndex: "name",
            render: (name: string) => {
                return <Link to={`/git/${name}`}>{name}</Link>;
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
                        <Button type={"link"} onClick={() => navigate(`/git/${record.name}/edit`)}>
                            Edit
                        </Button>
                        <Popconfirm
                            title={t`Are you sure to delete this git repository?`}
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
                    <Button icon={<PlusOutlined />} type={"primary"} onClick={() => navigate("/git/create")}>
                        {t`Create`}
                    </Button>
                </Space>
            }
        >
            <Table<ApiV1GitGitRepo>
                rowKey={"name"}
                loading={loading}
                columns={columns}
                dataSource={data?.list}
                pagination={pagination}
            />
        </Card>
    );
}