import { Button, Descriptions, Spin } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import { useRequest } from "ahooks";
import useApp from "antd/es/app/useApp";
import { buildSettingsApi } from "@/api";
import { EditOutlined } from "@ant-design/icons";

export default function DockerfileDescription() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const { data, loading } = useRequest(buildSettingsApi.buildSettingsServiceGetBuildSettings.bind(buildSettingsApi), {
        ready: !!name,
        defaultParams: [{ name: name as string }],
        onError: (error) => {
            notification.error({
                message: t`failed to get build settings`,
                description: error.message
            });
        }
    });
    const items: DescriptionsItemType[] = [
        {
            label: "",
            children: <div className={"w-full min-h-96 border p-2 border-gray-300"}>{data?.dockerfile}</div>
        }
    ];
    return (
        <Spin spinning={loading}>
            <Descriptions
                title={"Dockerfile"}
                extra={
                    <Button
                        type={"primary"}
                        size={"middle"}
                        icon={<EditOutlined />}
                    >
                        {t`Edit`}
                    </Button>
                }
                items={items}
            />
        </Spin>
    );
}