import {useLingui} from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import {useRequest} from "ahooks";
import {buildSettingsApi} from "@/api";
import {useNavigate, useParams} from "react-router";
import {useForm} from "antd/es/form/Form";
import type {ApiV1BuildSettingsBuildSettings} from "@/generated";
import {Button, Form, Input, Space, Spin} from "antd";
import {SaveOutlined} from "@ant-design/icons";

export default function BuildSettingsEdit() {
    const {id} = useParams();
    const {t} = useLingui();
    const {notification} = useApp();
    const navigate = useNavigate();
    const next = `/app/${id}/build`

    const [form] = useForm<ApiV1BuildSettingsBuildSettings>();
    const {Item} = Form;

    const {loading} = useRequest(buildSettingsApi.buildSettingsServiceGetBuildSettings.bind(buildSettingsApi), {
        ready: !!id,
        defaultParams: [{applicationId: id as string}],
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

    const {run: updateBuildSettings, loading: updateLoading} = useRequest(
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
        if (id) {
            updateBuildSettings({
                applicationId: id,
                apiV1BuildSettingsBuildSettings: values
            });
        }
    };


    return (
        <Spin spinning={loading}>
            <Form
                form={form}
                onFinish={submitForm}
                labelWrap={true}
                labelCol={{
                    flex: "100px"
                }}
            >
                <Item label={t`Dockerfile`} name={"dockerfile"}>
                    <Input.TextArea rows={4}/>
                </Item>
            </Form>
            <Space className={"ml-25"} size={"large"}>
                <Button icon={<SaveOutlined/>} type="primary" onClick={() => form.submit()} loading={updateLoading}>
                    {t`Save`}
                </Button>
                <Button onClick={() => navigate(next)}>{t`Cancel`}</Button>
            </Space>
        </Spin>
    );
}