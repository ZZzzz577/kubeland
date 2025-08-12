import { useLingui } from "@lingui/react/macro";
import { Form, Input, Radio } from "antd";
import UploadKubeConfig from "@/views/cluster/modify/components/form/UploadKubeConfig.tsx";
import { type ApiV1ClusterCluster, ApiV1ClusterConnectionTypeEnum } from "@/generated";
import useFormInstance from "antd/es/form/hooks/useFormInstance";

export default function ConnectionForm() {
    const { t } = useLingui();
    const form = useFormInstance<ApiV1ClusterCluster>();
    const connectionType = Form.useWatch(["connection", "type"], form);
    const { Item } = Form;
    return (
        <>
            <UploadKubeConfig />
            <Item
                label={t`address`}
                name={["connection", "address"]}
                rules={[
                    {
                        required: true,
                        pattern: /^https:\/\/[a-zA-Z0-9-.]+(:[0-9]+)?$/,
                    },
                ]}
            >
                <Input />
            </Item>
            <Item label={t`CA`} name={["connection", "ca"]} rules={[{ required: true }]}>
                <Input.TextArea autoSize={{ minRows: 4, maxRows: 8 }} />
            </Item>
            <Item label={t`type`} name={["connection", "type"]} initialValue={2} rules={[{ required: true }]}>
                <Radio.Group
                    optionType={"button"}
                    options={[
                        { label: t`CERT`, value: ApiV1ClusterConnectionTypeEnum.TlsCert },
                        { label: t`TOKEN`, value: ApiV1ClusterConnectionTypeEnum.TlsToken },
                    ]}
                />
            </Item>
            {connectionType === ApiV1ClusterConnectionTypeEnum.TlsCert && (
                <>
                    <Item label={t`cert`} name={["connection", "cert"]} rules={[{ required: true }]}>
                        <Input.TextArea autoSize={{ minRows: 4, maxRows: 8 }} />
                    </Item>
                    <Item label={t`key`} name={["connection", "key"]} rules={[{ required: true }]}>
                        <Input.TextArea autoSize={{ minRows: 4, maxRows: 8 }} />
                    </Item>
                </>
            )}
            {connectionType === ApiV1ClusterConnectionTypeEnum.TlsToken && (
                <Item label={"token"} name={["connection", "token"]} rules={[{ required: true }]}>
                    <Input.TextArea autoSize={{ minRows: 4, maxRows: 8 }} />
                </Item>
            )}
        </>
    );
}
