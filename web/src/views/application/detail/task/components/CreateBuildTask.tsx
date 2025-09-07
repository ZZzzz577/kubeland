import { Button } from "antd";
import { useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import { PlusCircleOutlined } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { buildTaskApi } from "@/api";
import useApp from "antd/es/app/useApp";

export default function CreateBuildTask(props: { refresh: () => void }) {
    const { refresh } = props;
    const { name } = useParams();
    const { t } = useLingui();
    const { notification } = useApp();
    const { run, loading } = useRequest(buildTaskApi.buildTaskServiceCreate.bind(buildTaskApi), {
        manual: true,
        onSuccess: () => {
            notification.success({ message: t`Create build task success` });
            refresh()
        },
        onError: (e) =>
            notification.error({
                message: t`Create build task failed`,
                description: e.message,
            }),
    });
    const onClick = () => {
        if (name) {
            run({ name: name, apiV1ApplicationIdentityRequest: { name: name } });
        }
    };
    return (
        <Button
            loading={loading}
            onClick={onClick}
            type={"primary"}
            size={"middle"}
            icon={<PlusCircleOutlined />}
        >{t`Create`}</Button>
    );
}
