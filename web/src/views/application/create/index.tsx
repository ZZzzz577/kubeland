import { Card } from "antd";
import BackButton from "@/components/back/BackButton.tsx";
import { useLingui } from "@lingui/react/macro";
import ApplicationCreateForm from "@/views/application/create/components/ApplicationCreateForm.tsx";

export default function ApplicationCreate() {
    const { t } = useLingui();
    return (
        <Card title={t`Create a new application`} extra={<BackButton />}>
            <ApplicationCreateForm />
        </Card>
    );
}
