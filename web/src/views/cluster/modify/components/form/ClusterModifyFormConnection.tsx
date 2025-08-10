import { useLingui } from "@lingui/react/macro";
import { Button, Flex, Form, Input } from "antd";

export default function ClusterModifyFormConnection() {
    const { t } = useLingui();
    const { Item } = Form;
    return (
        <>
            <Flex className={"!mb-5"} align={"center"} gap={"middle"}>
                <div className={"text-base font-medium"}>{t`connection settings`}</div>
                <Button type={"link"} size={"small"}>{t`upload kube/config`}</Button>
            </Flex>
            <Item
                label={t`address`}
                name={"address"}
                rules={[
                    {
                        required: true,
                        pattern: /^https:\/\/[a-zA-Z0-9-.]+(:[0-9]+)?$/,
                    },
                ]}
            >
                <Input />
            </Item>
        </>
    );
}
