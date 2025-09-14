import { Button, Card, Form, Input, Space } from "antd";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { imageApi } from "@/api";
import { useNavigate } from "react-router";
import { useForm } from "antd/es/form/Form";
import type { ApiV1ImageImageRepo } from "@/generated";
import { SaveOutlined } from "@ant-design/icons";

export default function ImageRepositoryCreate() {
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { loading, run } = useRequest(imageApi.imageServiceCreateImageRepo.bind(imageApi), {
        manual: true,
        onSuccess: () => {
            notification.success({
                message: t`create image repository success`,
            });
            setTimeout(() => navigate("/image"), 500);
        },
        onError: (error) => {
            notification.error({
                message: t`failed to create image repository`,
                description: error.message,
            });
        },
    });

    const [form] = useForm<ApiV1ImageImageRepo>();
    const { Item } = Form;
    const onFinish = (values: ApiV1ImageImageRepo) => {
        run({
            apiV1ImageImageRepo: values,
        });
    };
    return (
        <Card title={t`Create image repository`}>
            <Form<ApiV1ImageImageRepo>
                form={form}
                onFinish={onFinish}
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
                <Item label={t`Description`} name={"description"}>
                    <Input.TextArea rows={4} />
                </Item>
                <Item
                    label={t`URL`}
                    name={"url"}
                    rules={[
                        {
                            required: true,
                            max: 256,
                        },
                    ]}
                >
                    <Input />
                </Item>
                <Item label={t`Username`} name={"username"}>
                    <Input />
                </Item>
                <Item label={t`Password`} name={"password"}>
                    <Input />
                </Item>
                <Space className={"ml-25"} size={"large"}>
                    <Button icon={<SaveOutlined />} type="primary" htmlType="submit" loading={loading}>
                        {t`Save`}
                    </Button>
                    <Button onClick={() => navigate(-1)}>{t`Cancel`}</Button>
                </Space>
            </Form>
        </Card>
    );
}
