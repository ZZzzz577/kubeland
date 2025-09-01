import { useLingui } from "@lingui/react/macro";
import useApp from "antd/es/app/useApp";
import { useRequest } from "ahooks";
import { applicationApi } from "@/api";
import { Button, Popconfirm } from "antd";
import type { BaseButtonProps } from "antd/es/button/button";
import { useParams } from "react-router";

export default function DeleteButton(props: BaseButtonProps) {
    const { id } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const { run, loading } = useRequest(applicationApi.applicationServiceDeleteApplication.bind(applicationApi), {
        manual: true,
        onSuccess: () => {
            notification.success({
                message: t`delete application success`,
            });
        },
        onError: (error) => {
            notification.error({
                message: t`delete application failed`,
                description: error.message,
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
            <Button {...props} danger loading={loading}>
                Delete
            </Button>
        </Popconfirm>
    );
}