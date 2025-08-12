import { Button, Card } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useNavigate, useParams } from "react-router";
import ClusterModifyForm from "@/views/cluster/modify/components/form/ClusterModifyForm.tsx";
import { DoubleLeftOutlined } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { clusterApi } from "@/api";

export default function ClusterModify() {
    const { t } = useLingui();
    const navigate = useNavigate();
    const { id } = useParams();
    const isUpdate = !!id;
    const title = isUpdate ? t`create cluster` : t`modify cluster`;

    const { data } = useRequest(() => (id ? clusterApi.clusterServiceGetCluster({ id }) : Promise.resolve(undefined)), {
        ready: isUpdate,
        refreshDeps: [id],
    });

    return (
        <Card
            title={title}
            extra={
                <Button
                    size={"large"}
                    icon={<DoubleLeftOutlined />}
                    type={"link"}
                    onClick={() => navigate("/cluster")}
                >{t`back`}</Button>
            }
        >
            <ClusterModifyForm cluster={data} />
        </Card>
    );
}
