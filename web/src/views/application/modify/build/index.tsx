import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { buildSettingsApi } from "@/api";
import { useNavigate, useParams } from "react-router";
import { useForm } from "antd/es/form/Form";
import type { ApiV1BuildSettingsBuildSettings } from "@/generated";
import { Button, Card, Form, Input, Space, Spin } from "antd";
import { SaveOutlined } from "@ant-design/icons";
import type { ReactNode } from "react";
import GitEdit from "@/views/application/modify/build/components/GitEdit.tsx";
import ImageEdit from "@/views/application/modify/build/components/ImageEdit.tsx";

export const Title = (props: {
    children: ReactNode
}) => {
    return <div className={"text-base font-semibold  mb-6"}>{props.children}</div>;
};

export default function BuildSettingsEdit() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const next = `/app/${name}/build`;

    const [form] = useForm<ApiV1BuildSettingsBuildSettings>();
    const { Item } = Form;

    const { loading } = useRequest(buildSettingsApi.buildSettingsServiceGetBuildSettings.bind(buildSettingsApi), {
        ready: !!name,
        defaultParams: [{ name: name as string }],
        onSuccess: (data) => {
            form.setFieldsValue(data);
        },
        onError: (error) => {
            notification.error({
                message: t`failed to get build settings`,
                description: error.message
            });
        }
    });

    const { run: updateBuildSettings, loading: updateLoading } = useRequest(
        buildSettingsApi.buildSettingsServiceApplyBuildSettings.bind(buildSettingsApi),
        {
            manual: true,
            onSuccess: () => {
                notification.success({
                    message: t`update building settings success`
                });
                setTimeout(() => navigate(next), 500);
            },
            onError: (error) => {
                notification.error({
                    message: t`failed to update building settings`,
                    description: error.message
                });
            }
        }
    );
    const submitForm = (values: ApiV1BuildSettingsBuildSettings) => {
        if (name) {
            updateBuildSettings({
                name: name,
                apiV1BuildSettingsBuildSettings: values
            });
        }
    };


    return (
        <Card>
            <Spin spinning={loading}>
                <Form
                    form={form}
                    onFinish={submitForm}
                    labelWrap={true}
                    labelAlign={"left"}
                    labelCol={{
                        flex: "100px"
                    }}
                >

                    <GitEdit />
                    <ImageEdit />

                    <Title>{t`Dockerfile`}</Title>
                    <Item name={"dockerfile"} rules={[{ required: true }]}>
                        <Input.TextArea rows={8} />
                    </Item>
                </Form>
                <Space className={"ml-25"} size={"large"}>
                    <Button icon={<SaveOutlined />} type="primary" onClick={() => form.submit()}
                            loading={updateLoading}>
                        {t`Save`}
                    </Button>
                    <Button onClick={() => navigate(next)}>{t`Cancel`}</Button>
                </Space>
            </Spin>
        </Card>
    );
}