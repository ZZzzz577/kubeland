import { Descriptions } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { useLingui } from "@lingui/react/macro";
import type { ApiV1BuildSettingsBuildSettingsGitSettings } from "@/generated";

export default function GitDescription(props: {
    git?: ApiV1BuildSettingsBuildSettingsGitSettings
}) {
    const { git } = props;
    const { t } = useLingui();
    const items: DescriptionsItemType[] = [
        {
            label: t`URL`,
            children: git?.url
        }
    ];
    return <Descriptions title={"Git repository settings"} bordered styles={{ label: { width: 150 } }} items={items} />;
}