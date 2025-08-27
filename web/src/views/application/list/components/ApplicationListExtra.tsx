import { Button, Space } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { useLingui } from "@lingui/react/macro";
import { useNavigate } from "react-router";

export default function ApplicationListExtra() {
    const { t } = useLingui();
    const navigate = useNavigate();
    return (
        <Space>
            <Button icon={<PlusOutlined />} type={"primary"} onClick={() => navigate("/app/create")}>
                {t`Create`}
            </Button>
        </Space>
    );
}