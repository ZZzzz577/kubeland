import { Descriptions } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { useLingui } from "@lingui/react/macro";
import type { ApiV1BuildSettingsBuildSettingsGitSettings } from "@/generated";
import { Link } from "react-router";

export default function GitDescription(props: {
    git?: ApiV1BuildSettingsBuildSettingsGitSettings,
}) {
    const { git } = props;
    const { t } = useLingui();

    const items: DescriptionsItemType[] = [
        {
            label: t`URL`,
            children: <Link to={git?.url ?? ""} target={"_blank"}>{git?.url}</Link>
        }
    ];

    return (
        <Descriptions
            title={t`Git repository settings`}
            bordered
            column={1}
            styles={{ label: { width: 150 } }}
            items={items} />

    );
}