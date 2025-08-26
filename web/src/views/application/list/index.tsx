import { Card } from "antd";
import ApplicationListTable from "@/views/application/list/components/ApplicationListTable.tsx";
import ApplicationListExtra from "@/views/application/list/components/ApplicationListExtra.tsx";
import { useLingui } from "@lingui/react/macro";

export default function ApplicationList() {
    const { t } = useLingui();
    return (
        <Card title={t`Application list`} extra={<ApplicationListExtra />}>
            <ApplicationListTable />
        </Card>
    );
}
