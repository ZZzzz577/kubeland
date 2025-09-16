import { Button, Card, Space } from "antd";
import DockerfileDescription from "@/views/application/detail/build/components/DockerfileDescription.tsx";
import GitDescription from "@/views/application/detail/build/components/GitDescription.tsx";
import ImageDescription from "@/views/application/detail/build/components/ImageDescription.tsx";
import { useNavigate, useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { buildSettingsApi } from "@/api";
import { EditOutlined } from "@ant-design/icons";


export default function BuildSettings() {
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const navigate = useNavigate();
    const { data, loading } = useRequest(buildSettingsApi.buildSettingsServiceGetBuildSettings.bind(buildSettingsApi), {
        ready: !!name,
        defaultParams: [{ name: name as string }],
        onError: (error) => {
            notification.error({
                message: t`failed to get build settings`,
                description: error.message
            });
        }
    });

    return (
        <Card
            loading={loading}
            className={"!rounded-t-none"}
            variant="borderless"
            title={t`Build settings`}
            extra={<Button
                type={"primary"}
                size={"middle"}
                icon={<EditOutlined />}
                onClick={() => navigate(`/app/${name}/edit/build`)}
            >
                {t`Edit`}
            </Button>}
        >
            <Space direction={"vertical"} size={"large"}>
                <GitDescription git={data?.git} />
                <ImageDescription image={data?.image}/>
                <DockerfileDescription dockerfile={data?.dockerfile} />
            </Space>
        </Card>
    );

}