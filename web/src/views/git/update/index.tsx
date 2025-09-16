import { Button, Card, Form, Input, Space } from "antd";
import { useNavigate, useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { gitApi } from "@/api";
import { useForm } from "antd/es/form/Form";
import type { ApiV1GitGitRepo } from "@/generated";
import { SaveOutlined } from "@ant-design/icons";

export default function GitRepoUpdate() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();

    const [form] = useForm<ApiV1GitGitRepo>();
    const { Item } = Form;

    const { data, loading: getLoading } = useRequest(gitApi.gitServiceGetGitRepo.bind(gitApi), {
        defaultParams: [{ name: name as string }],
        ready: !!name,
        onSuccess: (data) => {
            form.setFieldsValue(data);
        },
        onError: (e) => {
            notification.error({
                message: t`failed to get git repository detail`,
                description: e.message,
            });
        },
    });

    const { run, loading } = useRequest(gitApi.gitServiceUpdateGitRepo.bind(gitApi), {
        manual: true,
        onSuccess: () => {
            notification.success({
                message: t`update git repository success`,
            });
            setTimeout(() => navigate(`/git/${name}`), 500);
        },
        onError: (e) => {
            notification.error({
                message: t`failed to update git repository`,
                description: e.message,
            });
        },
    });

    const onFinish = (values: ApiV1GitGitRepo) => {
        if (name) {
            run({
                name: name,
                apiV1GitGitRepo: values,
            });
        }
    };
    return (
        <Card loading={getLoading} title={name}>
            <Form<ApiV1GitGitRepo>
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
                <Item label={t`Token`} name={"token"}>
                    <Input.TextArea rows={4} />
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