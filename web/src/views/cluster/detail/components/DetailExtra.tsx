import { Button, Space } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useNavigate } from "react-router";
import BackButton from "@/components/back/BackButton.tsx";
import { EditOutlined } from "@ant-design/icons";

export default function DetailExtra(props: { id?: string }) {
    const { t } = useLingui();
    const navigate = useNavigate();
    const { id } = props;
    return (
        <Space>
            <Button
                icon={<EditOutlined />}
                type={"primary"}
                onClick={() => navigate(`/cluster/${id}/edit`)}
            >{t`Edit`}</Button>
            <BackButton />
        </Space>
    );
}
