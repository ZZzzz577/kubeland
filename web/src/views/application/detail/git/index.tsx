import { Card, Select, Space, Spin, Table } from "antd";
import { useLingui } from "@lingui/react/macro";
import { usePagination, useRequest } from "ahooks";
import { gitApi } from "@/api";
import { useParams } from "react-router";
import useApp from "antd/es/app/useApp";
import { useState } from "react";
import type { ApiV1GitListCommitsResponseCommit } from "@/generated";

export default function GitRepo() {
    const { t } = useLingui();
    const { name } = useParams();
    const { notification } = useApp();
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
                columns={[
                    { title: t`SHA`, dataIndex: "sha", render: (value?: string) => value?.slice(0, 7) },
                    { title: t`Message`, dataIndex: "message", ellipsis: true },
                    {
                        title: t`Create time`,
                        dataIndex: "createdAt",
                        render: (value?: Date) => value?.toLocaleString(),
                    },
                ]}
                className={"mt-4"}
                rowKey={"sha"}
                loading={commitsLoading}
                pagination={pagination}
                dataSource={commitsData?.list}
            />
        </Card>
    );
}
