import type { ApiV1ApplicationApplication } from "@/generated";
import { useLingui } from "@lingui/react/macro";
import { App, Button, Form, Input, Select, Space, Spin } from "antd";
import { useNavigate } from "react-router";
import { useForm } from "antd/es/form/Form";
import { useEffect } from "react";
import { SaveOutlined } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { applicationApi, clusterApi } from "@/api";

export default function ApplicationForm(props: { app?: ApiV1ApplicationApplication }) {
    const { app } = props;
    const isUpdate = !!app;
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
                console.log(error);
                notification.error({
                    message: t`create application failed`,
                    description: error.message,
                });
            },
        },
    );

    const { run: updateApp, loading: updateLoading } = useRequest(
        applicationApi.applicationServiceUpdateApplication.bind(applicationApi),
        {
            manual: true,
            onSuccess: () => {
                notification.success({
                    message: t`update application success`,
                });
                setTimeout(() => navigate("/app"), 500);
            },
            onError: (error) => {
                notification.error({
                    message: t`update application failed`,
                    description: error.message,
                });
            },
        },
    );

    const [form] = useForm<ApiV1ApplicationApplication>();
    const { Item } = Form;
    useEffect(() => {
        if (app) {
            form.setFieldsValue(app);
        }
    }, [app, form]);
    const submitForm = (values: ApiV1ApplicationApplication) => {
        if (isUpdate) {
            if (app?.id) {
                updateApp({
                    id: app.id,
                    apiV1ApplicationApplication: values,
                });
            }
        } else {
            createApp({
                apiV1ApplicationApplication: values,
            });
        }
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
                {isUpdate ? <div>{app?.name}</div> : <Input />}
            </Item>
            <Item label={t`Cluster`} name={"clusterId"} rules={[{ required: true }]}>
                {isUpdate ? (
                    <Spin spinning={clusterLoading}>
                        {clusterList?.items?.find((item) => item.id === app.clusterId)?.name || ""}
                    </Spin>
                ) : (
                    <Select
                        loading={clusterLoading}
                        options={clusterList?.items?.map((item) => ({
                            label: item.name,
                            value: item.id,
                        }))}
                    />
                )}
            </Item>
            <Item label={t`Description`} name={"description"}>
                <Input.TextArea rows={4} />
            </Item>

            <Space className={"ml-25"} size={"large"}>
                <Button
                    icon={<SaveOutlined />}
                    type="primary"
                    htmlType="submit"
                    loading={isUpdate ? updateLoading : createLoading}
                >
                    {t`Save`}
                </Button>
                <Button onClick={() => navigate(-1)}>{t`Cancel`}</Button>
            </Space>
        </Form>
    );
}