import { useLingui } from "@lingui/react/macro";
import { useForm } from "antd/es/form/Form";
import type { ApiV1ClusterCluster } from "@/generated";
import { Button, Divider, Form, notification, Space } from "antd";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import BasicInfoForm from "@/views/cluster/modify/components/form/BasicInfoForm.tsx";
import ConnectionForm from "@/views/cluster/modify/components/form/ConnectionForm.tsx";
import TestConnection from "@/views/cluster/modify/components/form/TestConnection.tsx";
import { useNavigate } from "react-router";
import { useEffect } from "react";
import { SaveOutlined } from "@ant-design/icons";

export default function ClusterModifyForm(props: { cluster?: ApiV1ClusterCluster }) {
    const { cluster } = props;
    const isUpdate = !!cluster;
    const { t } = useLingui();
    const [notify, notifyContext] = notification.useNotification();
    const navigate = useNavigate();

    const [form] = useForm<ApiV1ClusterCluster>();
    useEffect(() => {
        if (cluster) {
            form.setFieldsValue(cluster);
        }
    }, [cluster, form]);

    const { run: createCluster, loading: createLoading } = useRequest(
        clusterApi.clusterServiceCreateCluster.bind(clusterApi),
        {
            manual: true,
            onSuccess: () => {
                notify.success({
                    message: t`create cluster success`,
                });
                setTimeout(() => navigate("/cluster"), 500);
            },
            onError: (error) => {
                console.log(error);
                notify.error({
                    message: t`create cluster failed`,
                    description: error.message,
                });
            },
        },
    );

    const { run: updateCluster, loading: updateLoading } = useRequest(
        clusterApi.clusterServiceUpdateCluster.bind(clusterApi),
        {
            manual: true,
            onSuccess: () => {
                notify.success({
                    message: t`update cluster success`,
                });
                setTimeout(() => navigate("/cluster"), 500);
            },
            onError: (error) => {
                console.log(error);
                notify.error({
                    message: t`update cluster failed`,
                    description: error.message,
                });
            },
        },
    );

    const submitForm = (values: ApiV1ClusterCluster) => {
        if (isUpdate) {
            if (cluster?.id) {
                updateCluster({
                    id: cluster.id,
                    apiV1ClusterCluster: values,
                });
            }
        } else {
            createCluster({
                apiV1ClusterCluster: values,
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
            {notifyContext}
            <BasicInfoForm />
            <Divider />
            <ConnectionForm />

            <Space className={"ml-25"} size={"large"}>
                <TestConnection />
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