import { useLingui } from "@lingui/react/macro";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { ApplicationCluster } from "@/views/application/list/components/ApplicationListTable.tsx";
import { Descriptions, Spin } from "antd";
import { useRequest } from "ahooks";
import { useParams } from "react-router";
import useApp from "antd/es/app/useApp";
import { applicationApi } from "@/api";

export default function BasicInfo() {
    const { t } = useLingui();
    const { name } = useParams();
    const { notification } = useApp();
    const { data, loading } = useRequest(applicationApi.applicationServiceGetApplication.bind(applicationApi), {
        ready: !!name,
        refreshDeps: [name],
        defaultParams: [{ name: name as string }],
        onError: (e) => {
            notification.error({
                message: t`failed to get application detail`,
                description: e.message,
            });
        },
    });
    const items: DescriptionsItemType[] = [
        {
            label: t`Name`,
            span: "filled",
            children: data?.name,
        },
        {
            label: t`Cluster`,
            span: "filled",
            children: <ApplicationCluster id={data?.clusterId} />,
        },
        {
            label: t`Description`,
            span: "filled",
            children: data?.description,
        },
        {
            label: t`Create time`,
            children: data?.createdAt?.toLocaleString(),
        },
        {
            label: t`Update time`,
            children: data?.updateAt?.toLocaleString(),
        },
    ];

    return (
        <Spin spinning={loading}>
            <Descriptions
                title={t`Basic info`}
                styles={{ label: { width: 150 } }}
                column={2}
                bordered
                items={items}
            />
        </Spin>
    );
}