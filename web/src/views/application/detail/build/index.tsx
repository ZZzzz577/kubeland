import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { buildSettingsApi } from "@/api";
import { Descriptions, Spin } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { useParams } from "react-router";

export default function BuildSettings() {
    const { id: appId } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const { data, loading } = useRequest(buildSettingsApi.buildSettingsServiceGetBuildSettings.bind(buildSettingsApi), {
        ready: !!appId,
        defaultParams: [{ applicationId: appId as string }],
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
            <Descriptions title={<div className={"text-sm"}>Dockerfile</div>} items={items} />
        </Spin>
    );

}