import { Descriptions } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { useLingui } from "@lingui/react/macro";
import type { ApiV1BuildSettingsBuildSettingsImageSettings } from "@/generated";

export default function ImageDescription(props: {
    image?: ApiV1BuildSettingsBuildSettingsImageSettings
}) {
    const { image } = props;
    const { t } = useLingui();
    const items: DescriptionsItemType[] = [
        {
            label: t`URL`,
            children: image?.url
        }
    ];
    return (
        <Descriptions
            title={"Image repository settings"}
            bordered
            styles={{ label: { width: 150 } }}
            items={items} />
    );
}