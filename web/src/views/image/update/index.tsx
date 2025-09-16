import { Button, Card, Form, Input, Space } from "antd";
import { useNavigate, useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { imageApi } from "@/api";
import { useForm } from "antd/es/form/Form";
import type { ApiV1ImageImageRepo } from "@/generated";
import { SaveOutlined } from "@ant-design/icons";

export default function ImageRepoUpdate() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();

    const [form] = useForm<ApiV1ImageImageRepo>();
    const { Item } = Form;

    const { data, loading: getLoading } = useRequest(imageApi.imageServiceGetImageRepo.bind(imageApi), {
        defaultParams: [{ name: name as string }],
        ready: !!name,
        onSuccess: (data) => {
            form.setFieldsValue(data);
        },
        onError: (e) => {
            notification.error({
                message: t`failed to get image repository detail`,
                description: e.message,
            });
        },
    });

    const { run, loading } = useRequest(imageApi.imageServiceUpdateImageRepo.bind(imageApi), {
        manual: true,
        onSuccess: () => {
            notification.success({
                message: t`update image repository success`,
            });
            setTimeout(() => navigate(`/image/${name}`), 500);
        },
        onError: (e) => {
            notification.error({
                message: t`failed to update image repository`,
                description: e.message,
            });
        },
    });

    const onFinish = (values: ApiV1ImageImageRepo) => {
        if (name) {
            run({
                name: name,
                apiV1ImageImageRepo: values,
            });
        }
    };
    return (
        <Card loading={getLoading} title={name}>
            <Form<ApiV1ImageImageRepo>
                form={form}
                onFinish={onFinish}
                labelWrap={true}
                labelCol={{
                    flex: "100px",
                }}
            >
                <Item label={t`Name`} required>
                    {data?.name}
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