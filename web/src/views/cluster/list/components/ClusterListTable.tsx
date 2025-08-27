import { usePagination, useRequest } from "ahooks";
import { clusterApi } from "@/api";
import { Spin, Table, Tag, Tooltip } from "antd";
import ClusterListTableColumns from "@/views/cluster/list/components/ClusterListTableColumns.tsx";
import type { ApiV1ClusterConnection } from "@/generated";
import { useState } from "react";
import useApp from "antd/es/app/useApp";
import { useLingui } from "@lingui/react/macro";

export default function ClusterListTable() {
    const { t } = useLingui();
    const { notification } = useApp();
    const { loading, pagination, data } = usePagination(
        async ({ current, pageSize }) => {
            const res = await clusterApi.clusterServiceListClusters({ pageCurrent: current, pageSize });
            return {
                current: res.pagination?.current ?? 0,
                total: res.pagination?.total ?? 0,
                list: res.items ?? [],
            };
        },
        {
            onError: (e) => {
                notification.error({
                    message: t`failed to get cluster list`,
                    description: e.message,
                });
            },
        },
    );

    const columns = ClusterListTableColumns();
    return <Table rowKey={"id"} loading={loading} pagination={pagination} columns={columns} dataSource={data?.list} />;
}

interface ClusterStatusDetail {
    tag: "unknown" | "success" | "failed";
    description: string;
}

export function ClusterStatus(props: { conn?: ApiV1ClusterConnection }) {
    const { conn } = props;
    const [status, setStatus] = useState<ClusterStatusDetail>({
        tag: "unknown",
        description: "",
    });

    const { loading } = useRequest(clusterApi.clusterServiceTestConnection.bind(clusterApi), {
        ready: !!conn,
        refreshDeps: [conn],
        defaultParams: [{ apiV1ClusterConnection: conn as ApiV1ClusterConnection }],
        onSuccess(data) {
            if (data.version) {
                setStatus({
                    tag: "success",
                    description: data.version,
                });
            }
        },
        onError(error) {
            setStatus({
                tag: "failed",
                description: error.message || "Connection failed",
            });
        },
    });

    const tooltip = () => {
        switch (status.tag) {
            case "success":
                return `Connection successful. Version is ${status.description}`;
            case "failed":
                return `Connection failed. ${status.description}`;
            default:
                return `Connection status is unknown`;
        }
    };

    const renderTag = () => {
        switch (status.tag) {
            case "success":
                return <Tag color="green">success</Tag>;
            case "failed":
                return <Tag color="red">failed</Tag>;
            default:
                return <Tag>unknown</Tag>;
        }
    };

    return (
        <Spin spinning={loading}>
            <Tooltip title={tooltip()} className={"flex"}>
                {renderTag()}
                <div className={"w-25 truncate"}>{status.description}</div>
            </Tooltip>
        </Spin>
    );
}

interface OperatorStatusDetail {
    tag: "unknown" | "success" | "failed";
    description?: string;
}

export function OperatorStatus(props: { conn?: ApiV1ClusterConnection }) {
    const { conn } = props;
    const [status, setStatus] = useState<OperatorStatusDetail>({
        tag: "unknown",
        description: "",
    });

    const { loading } = useRequest(clusterApi.clusterServiceTestOperator.bind(clusterApi), {
        ready: !!conn,
        refreshDeps: [conn],
        defaultParams: [{ apiV1ClusterConnection: conn as ApiV1ClusterConnection }],
        onSuccess() {
            setStatus({
                tag: "success",
            });
        },
        onError() {
            setStatus({
                tag: "failed",
            });
        },
    });

    const renderTag = () => {
        switch (status.tag) {
            case "success":
                return <Tag color="green">success</Tag>;
            case "failed":
                return <Tag color="red">failed</Tag>;
            default:
                return <Tag>unknown</Tag>;
        }
    };

    return (
        <Spin spinning={loading}>
            <Tooltip title={status.description} className={"flex"}>
                {renderTag()}
                <div className={"w-25 truncate"}>{status.description}</div>
            </Tooltip>
        </Spin>
    );
}
