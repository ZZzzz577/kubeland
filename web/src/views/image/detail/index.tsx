import { useNavigate, useParams } from "react-router";
import useApp from "antd/es/app/useApp";
import { Button, Card, Descriptions, Space } from "antd";
import { useRequest } from "ahooks";
import { imageApi } from "@/api";
import { useLingui } from "@lingui/react/macro";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { EditOutlined } from "@ant-design/icons";

export default function ImageRepoDetail() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { loading, data } = useRequest(imageApi.imageServiceGetImageRepo.bind(imageApi), {
        defaultParams: [{ name: name as string }],
        ready: !!name,
        onError: (e) => {
            notification.error({
                message: t`failed to get image repository detail`,
                description: e.message,
            });
        },
    });
    const items: DescriptionsItemType[] = [
        {
            label: t`Name`,
            children: data?.name,
            span: "filled",
        },
        {
            label: t`Description`,
            children: data?.description,
            span: "filled",
        },
        {
            label: t`URL`,
            children: data?.url,
            span: "filled",
        },
        {
            label: t`Username`,
            children: data?.username,
            span: "filled",
        },
        {
            label: t`Password`,
            children: data?.password,
            span: "filled",
        },

        {
            label: t`Create Time`,
            children: data?.createdAt?.toLocaleDateString(),
        },
        {
            label: t`Update Time`,
            children: data?.updatedAt?.toLocaleDateString(),
        },
    ];
    return (
        <Card loading={loading}>
            <Descriptions
                title={name}
                column={2}
                bordered
                styles={{
                    label: { width: 150 },
                }}
                items={items}
                extra={
                    <Space>
                        <Button
                            icon={<EditOutlined />}
                            type={"primary"}
                            onClick={() => navigate(`/image/${name}/edit`)}
                        >
                            Edit
                        </Button>
                    </Space>
                }
            />
        </Card>
    );
}