import type { ApiV1ClusterCluster } from "@/generated";
import { useLingui } from "@lingui/react/macro";
import { Descriptions, type DescriptionsProps } from "antd";

export default function BasicInfoDescription(props: { cluster?: ApiV1ClusterCluster }) {
    const { t } = useLingui();
    const { cluster } = props;
    const items: DescriptionsProps["items"] = [
        {
            label: t`Name`,
            children: cluster?.name,
            span: 2,
        },
        {
            label: t`Description`,
            children: cluster?.description,
            span: 2,
        },
        {
            label: t`Create time`,
            children: cluster?.createdAt?.toLocaleString(),
            span: 1,
        },
        {
            label: t`Update time`,
            children: cluster?.updatedAt?.toLocaleString(),
            span: 1,
        },
    ];
    return (
        <Descriptions
            title={t`Basic info`}
            column={2}
            bordered
            items={items}
            styles={{
                label: { width: 120, padding: 10 },
                root: { border: "none", borderRadius: 0 },
            }}
        />
    );
}