import { Card, Space } from "antd";
import { Trans, useLingui } from "@lingui/react/macro";
import type { CardTabListType } from "antd/es/card/Card";
import { type ReactNode, useMemo, useState } from "react";
import { BuildOutlined, MenuOutlined, SettingOutlined } from "@ant-design/icons";
import BasicInfo from "@/views/application/detail/components/BasicInfo.tsx";
import DetailExtra from "@/views/application/detail/components/DetailExtra.tsx";
import { useRequest } from "ahooks";
import { applicationApi } from "@/api";
import { useParams } from "react-router";
import useApp from "antd/es/app/useApp";

const tabList: CardTabListType[] = [
    {
        key: "basicInfo",
        label: (
            <Space>
                <MenuOutlined />
                <Trans>Basic info</Trans>
            </Space>
        ),
    },
    {
        key: "buildSettings",
        label: (
            <Space>
                <SettingOutlined />
                <Trans>Build settings</Trans>
            </Space>
        ),
    },
    {
        key: "buildTasks",
        label: (
            <Space>
                <BuildOutlined />
                <Trans>Build tasks</Trans>
            </Space>
        ),
    },
];

export default function ApplicationDetail() {
    const { t } = useLingui();
    const { id } = useParams();
    const { notification } = useApp();

    const [activeTab, setActiveTab] = useState<string>("basicInfo");
    const onTabChange = (key: string) => {
        setActiveTab(key);
    };

    const { data, loading } = useRequest(applicationApi.applicationServiceGetApplication.bind(applicationApi), {
        ready: !!id,
        refreshDeps: [id],
        defaultParams: [{ id: id as string }],
        onError: (e) => {
            notification.error({
                message: t`failed to get application detail`,
                description: e.message,
            });
        },
    });

    const tabContent = useMemo((): ReactNode => {
        switch (activeTab) {
            case "basicInfo":
                return <BasicInfo info={data} />;
            default:
                return <></>;
        }
    }, [activeTab, data]);

    return (
        <Card
            loading={loading}
            title={<div className={"text-xl"}>{data?.name}</div>}
            tabList={tabList}
            activeTabKey={activeTab}
            onTabChange={onTabChange}
            tabProps={{ size: "middle" }}
            extra={<DetailExtra />}
        >
            {tabContent}
        </Card>
    );
}
