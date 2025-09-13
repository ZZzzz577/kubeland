import { Button, Descriptions, Spin } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";
import { useLingui } from "@lingui/react/macro";
import { useRequest } from "ahooks";
import { gitApi } from "@/api";
import useApp from "antd/es/app/useApp";
import { useParams } from "react-router";
import { EditOutlined } from "@ant-design/icons";
import { useState } from "react";
import GitForm from "@/views/application/detail/build/components/GitForm.tsx";

export default function GitDescription() {
    const {t} = useLingui();
    const {name} = useParams();
    const {notification} = useApp();

    const [editMode, setEditMode] = useState(false);

    const {loading, data, refresh} = useRequest(gitApi.gitServiceGetGitSettings.bind(gitApi), {
        ready: !!name,
        defaultParams: [{name: name as string}],
        onError: (e) => {
            notification.error({
                message: t`failed to get git settings`,
                description: e.message
            });
        }
    });

    const finish = () => {
        setEditMode(false);
        refresh();
    };


    const items: DescriptionsItemType[] = [
        {
            label: t`URL`,
            children: data?.url
        },
        {
            label: "Token",
            children: data?.token
        }
    ];
    return (
        <Spin spinning={loading}>
            {editMode ? <GitForm git={data} finish={finish}/> :
                <Descriptions
                    title={t`Git repository settings`}
                    extra={
                        <Button
                            type={"primary"}
                            size={"middle"}
                            icon={<EditOutlined/>}
                            onClick={() => setEditMode(true)}
                        >
                            {t`Edit`}
                        </Button>
                    }
                    bordered
                    column={1}
                    styles={{label: {width: 150}}}
                    items={items}/>}
        </Spin>
    );
}