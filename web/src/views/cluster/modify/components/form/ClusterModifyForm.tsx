import { useLingui } from "@lingui/react/macro";
import { useForm } from "antd/es/form/Form";
import type { ApiV1ClusterCluster } from "@/generated";
import {Button, Form, Space} from "antd";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";
import type { NotificationInstance } from "antd/es/notification/interface";
import ClusterModifyFormBasic from "@/views/cluster/modify/components/form/ClusterModifyFormBasic.tsx";
import ClusterModifyFormConnection from "@/views/cluster/modify/components/form/ClusterModifyFormConnection.tsx";

export default function ClusterModifyForm(props: { notify: NotificationInstance }) {
    const { notify } = props;
    const { t } = useLingui();
    const [form] = useForm<ApiV1ClusterCluster>();
    const createCluster = clusterApi.clusterServiceCreateCluster.bind(clusterApi);
    const { run, loading } = useRequest(createCluster, {
        manual: true,
        onSuccess: () => {
            notify.success({
                message: t`create cluster success`,
            });
        },
        onError: (error) => {
            console.log(error);
            notify.error({
                message: t`create cluster failed`,
                description: error.message,
            });
        },
    });
    const submitForm = (values: ApiV1ClusterCluster) => {
        run({
            apiV1ClusterCluster: values,
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
            <ClusterModifyFormBasic />
            <ClusterModifyFormConnection />
            <Space>
                <Button type="primary" htmlType="submit" loading={loading}>
                    {t`create`}
                </Button>
            </Space>
        </Form>
    );
}
