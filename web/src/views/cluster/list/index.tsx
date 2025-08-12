import { Button, Card, Space } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useNavigate } from "react-router";
import { PlusOutlined } from "@ant-design/icons";
import ClusterListTable from "@/views/cluster/list/components/table/ClusterListTable.tsx";

export default function ClusterList() {
    const { t } = useLingui();
    const navigate = useNavigate();
    return (
        <Card
            title={t`cluster list`}
            extra={
                <Space>
                    <Button icon={<PlusOutlined />} type={"primary"} onClick={() => navigate("/cluster/create")}>
                        {t`create`}
                    </Button>
                </Space>
            }
        >
            <ClusterListTable />
        </Card>
    );
}
