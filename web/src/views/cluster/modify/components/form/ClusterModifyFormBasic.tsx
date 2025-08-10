import { useLingui } from "@lingui/react/macro";
import { Form, Input } from "antd";

export default function ClusterModifyFormBasic() {
    const { t } = useLingui();
    const { Item } = Form;
    return (
        <>
            <div className={"text-base font-medium mb-5"}>{t`basic info`}</div>
            <Item
                label={t`name`}
                name={"name"}
                rules={[
                    {
                        required: true,
                        max: 64,
                        pattern: /^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$/,
                    },
                ]}
            >
                <Input />
            </Item>
            <Item label={t`description`} name={"description"}>
                <Input.TextArea rows={4} />
            </Item>
        </>
    );
}
