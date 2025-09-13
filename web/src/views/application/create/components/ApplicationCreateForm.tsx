import type { ApiV1ApplicationApplication } from "@/generated";
import { useLingui } from "@lingui/react/macro";
import { App, Button, Form, Input, Select, Space } from "antd";
import { useNavigate } from "react-router";
import { useForm } from "antd/es/form/Form";
import { SaveOutlined } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { applicationApi, clusterApi } from "@/api";

export default function ApplicationCreateForm() {
    const { t } = useLingui();
    const { notification } = App.useApp();
    const navigate = useNavigate();

    // todo add list all clusters
    const { data: clusterList, loading: clusterLoading } = useRequest(
        clusterApi.clusterServiceListClusters.bind(clusterApi),
    );

    const { run: createApp, loading: createLoading } = useRequest(
        applicationApi.applicationServiceCreateApplication.bind(applicationApi),
        {
            manual: true,
            onSuccess: () => {
                notification.success({
                    message: t`create application success`,
                });
                setTimeout(() => navigate("/app"), 500);
            },
            onError: (error) => {
                notification.error({
                    message: t`create application failed`,
                    description: error.message,
                });
            },
        },
    );

    const [form] = useForm<ApiV1ApplicationApplication>();
    const { Item } = Form;

    const submitForm = (values: ApiV1ApplicationApplication) => {
        createApp({
            apiV1ApplicationApplication: values,
        });
    };

    return (
        <Form
            form={form}
            onFinish={submitForm}
            labelWrap={true}
            labelCol={{
                flex: "100px",
            }}
        >
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
            <Item label={t`Cluster`} name={"clusterId"} rules={[{ required: true }]}>
                <Select
                    loading={clusterLoading}
                    options={clusterList?.items?.map((item) => ({
                        label: item.name,
                        value: item.id,
                    }))}
                />
            </Item>
            <Item label={t`Description`} name={"description"}>
                <Input.TextArea rows={4} />
            </Item>

            <Space className={"ml-25"} size={"large"}>
                <Button icon={<SaveOutlined />} type="primary" htmlType="submit" loading={createLoading}>
                    {t`Save`}
                </Button>
                <Button onClick={() => navigate(-1)}>{t`Cancel`}</Button>
            </Space>
        </Form>
    );
}
