import { Button, Card, Select, Space, Spin, Table } from "antd";
import { useLingui } from "@lingui/react/macro";
import { usePagination, useRequest } from "ahooks";
import { buildTaskApi, gitApi } from "@/api";
import { Link, useNavigate, useParams } from "react-router";
import useApp from "antd/es/app/useApp";
import { useState } from "react";
import type { ApiV1GitListCommitsResponseCommit } from "@/generated";
import type { ColumnsType } from "antd/es/table";

export default function GitRepo() {
    const { t } = useLingui();
    const { name } = useParams();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { loading, data } = useRequest(gitApi.gitServiceGetGitSettings.bind(gitApi), {
        ready: !!name,
        defaultParams: [{ name: name as string }],
        onError: (e) => {
            notification.error({
                message: t`failed to get git settings`,
                description: e.message,
            });
        },
    });

    const [selectBranch, setSelectBranch] = useState("main");
    const { loading: branchesLoading, data: branches } = useRequest(gitApi.gitServiceListBranches.bind(gitApi), {
        ready: !!name,
        defaultParams: [{ name: name as string }],
        onError: (e) => {
            notification.error({
                message: t`failed to list branches`,
                description: e.message,
            });
        },
    });

    const {
        loading: commitsLoading,
        data: commitsData,
        pagination,
    } = usePagination(
        async ({ current, pageSize }) => {
            const res = await gitApi.gitServiceListCommits({
                pageCurrent: current,
                pageSize,
                name: name as string,
                branchName: selectBranch,
            });
            return {
                current: res.pagination?.current ?? 0,
                total: (res.pagination?.totalPage ?? 0) * pageSize,
                list: res.items ?? [],
            };
        },
        {
            ready: !!name,
            refreshDeps: [name, selectBranch],
            onError: (e) => {
                notification.error({
                    message: t`failed to list commits`,
                    description: e.message,
                });
            },
        },
    );

    const { run: buildTaskRun, loading: buildTaskLoading } = useRequest(
        buildTaskApi.buildTaskServiceCreate.bind(buildTaskApi),
        {
            manual: true,
            onSuccess: () => {
                notification.success({
                    message: t`build task created`,
                });
                setTimeout(() => navigate(`/app/${name}/task`), 500);
            },
            onError: (e) => {
                notification.error({
                    message: t`failed to create build task`,
                    description: e.message,
                });
            },
        },
    );
    const buildTask = (commit?: string) => {
        if (name && commit) {
            buildTaskRun({
                appName: name,
                apiV1BuildTaskCreateBuildTaskRequest: {
                    appName: name,
                    commit: commit,
                },
            });
        }
    };

    const columns: ColumnsType<ApiV1GitListCommitsResponseCommit> = [
        {
            title: t`SHA`,
            dataIndex: "sha",
            render: (value?: string) => (
                <Link target={"_blank"} to={`${data?.url ?? ""}/commit/${value}`}>
                    {value?.slice(0, 7)}
                </Link>
            ),
        },
        { title: t`Message`, dataIndex: "message", ellipsis: true },
        {
            title: t`Create time`,
            dataIndex: "createdAt",
            render: (value?: Date) => value?.toLocaleString(),
        },
        {
            title: t`Action`,
            dataIndex: "sha",
            render: (value?: string) => (
                <Space>
                    <Button
                        type={"link"}
                        loading={buildTaskLoading}
                        onClick={() => buildTask(value)}
                    >{t`Build`}</Button>
                </Space>
            ),
        },
    ];

    return (
        <Card
            className={"!rounded-t-none"}
            variant="borderless"
            extra={
                <Space size={"large"}>
                    <Space>
                        <span className={"text-gray-400"}>URL: </span>
                        <Spin spinning={loading}>{data?.url}</Spin>
                    </Space>
                    <Space>
                        <span className={"text-gray-400"}>Branch: </span>
                        <Spin spinning={loading}>
                            <Select
                                value={selectBranch}
                                onChange={(value) => setSelectBranch(value)}
                                className={"min-w-60"}
                                loading={branchesLoading}
                                options={branches?.items?.map((item) => ({ label: item, value: item }))}
                            />
                        </Spin>
                    </Space>
                </Space>
            }
        >
            <Table<ApiV1GitListCommitsResponseCommit>
                columns={columns}
                className={"mt-4"}
                rowKey={"sha"}
                loading={commitsLoading}
                pagination={pagination}
                dataSource={commitsData?.list}
            />
        </Card>
    );
}