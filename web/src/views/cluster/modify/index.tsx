import { Card, notification } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useParams } from "react-router";
import ClusterModifyForm from "@/views/cluster/modify/components/form/ClusterModifyForm.tsx";

export default function ClusterModify() {
    const { t } = useLingui();
    const { id } = useParams();
    const [notify, norityContext] = notification.useNotification();
    const isCreate = !!id;
    const title = isCreate ? t`modify cluster` : t`create cluster`;

    return (
        <Card title={title}>
            {norityContext}
            <ClusterModifyForm notify={notify} />
        </Card>
    );
}
