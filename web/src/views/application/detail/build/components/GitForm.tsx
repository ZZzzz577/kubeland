import { useLingui } from "@lingui/react/macro";
import { useRequest } from "ahooks";
import { gitApi } from "@/api";
import useApp from "antd/es/app/useApp";
import { Button, Form, Input, Space } from "antd";
import type { ApiV1GitGitSettings } from "@/generated";
import { useParams } from "react-router";
import { SaveOutlined } from "@ant-design/icons";

export default function GitForm(props: { git?: ApiV1GitGitSettings, finish: () => void }) {
    const { git, finish } = props;
    const { t } = useLingui();
    const { notification } = useApp();
    const { name } = useParams();
    const { loading, run } = useRequest(gitApi.gitServiceApplyGitSettings.bind(gitApi), {
        manual: true,
        onSuccess: () => {
            notification.success({ message: t`edit git settings success` });
            finish();
        },
        onError: (e) => {
            notification.error({
                message: t`failed to edit git settings`,
                description: e.message
            });
        }
    });
    const [form] = Form.useForm<ApiV1GitGitSettings>();
    const { Item } = Form;
    const onFinish = (values: ApiV1GitGitSettings) => {
        if (name) {
            run({
                name: name,
                apiV1GitGitSettings: values
            });
        }
    };
    return (
        <>
            <div className={"text-base font-semibold mb-5"}>{t`Git repository settings`}</div>
            <Form form={form}
                  onFinish={onFinish}
                  initialValues={git}
                  labelAlign={"left"}
                  labelCol={{
                      flex: "100px"
                  }}>
                <Item label={t`URL`} name={"url"} rules={[{ required: true }]}>
                    <Input />
                </Item>
                <Item label={t`Token`} name={"token"} rules={[{ required: true }]}>
                    <Input.TextArea rows={4}  />
                </Item>
                <Space className={"ml-25"} size={"large"}>
                    <Button icon={<SaveOutlined />} type="primary" onClick={() => form.submit()} loading={loading}>
                        {t`Save`}
                    </Button>
                    <Button onClick={finish}>{t`Cancel`}</Button>
                </Space>
            </Form>
        </>
    );
}