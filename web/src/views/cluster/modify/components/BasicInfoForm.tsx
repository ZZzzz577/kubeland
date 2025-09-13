import { useLingui } from "@lingui/react/macro";
import { Form, Input } from "antd";

export default function BasicInfoForm() {
    const { t } = useLingui();
    const { Item } = Form;
    return (
        <>
            <div className={"text-base font-medium mb-5"}>{t`Basic info`}</div>
            <Item
                label={t`Name`}
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
            <Item label={t`Description`} name={"description"}>
                <Input.TextArea rows={4} />
            </Item>
        </>
    );
}
