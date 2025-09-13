import { useLingui } from "@lingui/react/macro";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { ApplicationCluster } from "@/views/application/list/components/ApplicationListTable.tsx";
import { Button, Card, Descriptions, Space } from "antd";
import { useRequest } from "ahooks";
import { useNavigate, useParams } from "react-router";
import useApp from "antd/es/app/useApp";
import { applicationApi } from "@/api";
import { EditOutlined } from "@ant-design/icons";

export default function BasicInfo() {
    const { t } = useLingui();
    const { name } = useParams();
    const { notification } = useApp();
    const navigate = useNavigate();
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
        <Card
            title={t`Basic info`}
            className={"!rounded-t-none"}
            variant="borderless"
            loading={loading}
            extra={
                <Space>
                    <Button
                        type={"primary"}
                        size={"middle"}
                        icon={<EditOutlined />}
                        onClick={() => navigate(`/app/${name}/edit`)}
                    >{t`Edit`}</Button>
                </Space>
            }
        >
            <Descriptions styles={{ label: { width: 150 } }} column={2} bordered items={items} />
        </Card>
    );
}
