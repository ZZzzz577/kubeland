import type { ApiV1ApplicationApplication } from "@/generated";
import { useLingui } from "@lingui/react/macro";
import { App, Button, Form, Input, Space, Spin } from "antd";
import { useNavigate, useParams } from "react-router";
import { useForm } from "antd/es/form/Form";
import { SaveOutlined } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { applicationApi, clusterApi } from "@/api";

export default function BasicInfoEdit() {
    const { t } = useLingui();
    const { name } = useParams();
    const { notification } = App.useApp();
    const navigate = useNavigate();
    const next = `/app/${name}`;

    const [form] = useForm<ApiV1ApplicationApplication>();
    const { Item } = Form;

    const { data: app, loading } = useRequest(applicationApi.applicationServiceGetApplication.bind(applicationApi), {
        ready: !!name,
        refreshDeps: [name],
        defaultParams: [{ name: name as string }],
        onSuccess: (data) => {
            form.setFieldsValue(data);
        },
        onError: (e) => {
            notification.error({
                message: t`failed to get application detail`,
                description: e.message,
            });
        },
    });

    // todo add list all clusters
    const { data: clusterList, loading: clusterLoading } = useRequest(
        clusterApi.clusterServiceListClusters.bind(clusterApi),
    );

    const { run: updateApp, loading: updateLoading } = useRequest(
        applicationApi.applicationServiceUpdateApplication.bind(applicationApi),
        {
            manual: true,
            onSuccess: () => {
                notification.success({
                    message: t`update application success`,
                });
                setTimeout(() => navigate(next), 500);
            },
            onError: (error) => {
                notification.error({
                    message: t`update application failed`,
                    description: error.message,
                });
            },
        },
    );

    const submitForm = (values: ApiV1ApplicationApplication) => {
        if (name) {
            updateApp({
                name: name,
                apiV1ApplicationApplication: values,
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
                    flex: "100px",
                }}
            >
                <Item label={t`Name`} name={"name"} required>
                    <div>{app?.name}</div>
                </Item>
                <Item label={t`Cluster`} name={"clusterId"} required>
                    <Spin spinning={clusterLoading}>
                        {clusterList?.items?.find((item) => item.id === app?.clusterId)?.name || ""}
                    </Spin>
                </Item>
                <Item label={t`Description`} name={"description"}>
                    <Input.TextArea rows={4} />
                </Item>

                <Space className={"ml-25"} size={"large"}>
                    <Button icon={<SaveOutlined />} type="primary" htmlType="submit" loading={updateLoading}>
                        {t`Save`}
                    </Button>
                    <Button onClick={() => navigate(next)}>{t`Cancel`}</Button>
                </Space>
            </Form>
        </Spin>
    );
}