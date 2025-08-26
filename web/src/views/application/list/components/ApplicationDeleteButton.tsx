import { Button, Popconfirm } from "antd";
import { applicationApi } from "@/api";
import { useRequest } from "ahooks";
import useApp from "antd/es/app/useApp";
import { useLingui } from "@lingui/react/macro";

export default function ApplicationDeleteButton(props: { id?: string }) {
    const { id } = props;
    const { t } = useLingui();
    const { notification } = useApp();
    const { run, loading } = useRequest(applicationApi.applicationServiceDeleteApplication.bind(applicationApi), {
        manual: true,
        onSuccess: () => {
            notification.success({
                message: t`delete application success`,
            });
        },
    });
    const onClick = () => {
        if (id) {
            run({ id });
        }
    };
    return (
        <Popconfirm title={t`Are you sure to delete this application?`} onConfirm={onClick}>
            <Button type={"link"} danger loading={loading}>
                delete
            </Button>
        </Popconfirm>
    );
}