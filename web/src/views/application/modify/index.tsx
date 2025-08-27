import { useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import { useRequest } from "ahooks";
import { applicationApi } from "@/api";
import { Card } from "antd";
import BackButton from "@/components/back/BackButton.tsx";
import ApplicationModifyForm from "@/views/application/modify/components/ApplicationModifyForm.tsx";

export default function ApplicationModify() {
    const { t } = useLingui();
    const { id } = useParams();
    const isUpdate = !!id;
    const title = isUpdate ? t`Update cluster` : t`Create cluster`;

    const { data, loading } = useRequest(
        () => {
            return applicationApi.applicationServiceGetApplication({
                id: id as string,
            });
        },
        {
            ready: isUpdate,
            refreshDeps: [id],
        },
    );

    return (
        <Card title={title} loading={loading} extra={<BackButton />}>
            <ApplicationModifyForm app={data} />
        </Card>
    );
}