import { Card, Space } from "antd";
import { Trans } from "@lingui/react/macro";
import type { CardTabListType } from "antd/es/card/Card";
import { useMemo } from "react";
import { MenuOutlined, SettingOutlined } from "@ant-design/icons";
import DetailExtra from "@/views/application/detail/components/DetailExtra.tsx";
import { Outlet, useLocation, useNavigate, useParams } from "react-router";
import getActivePath from "@/views/application/common/utils/tab.ts";

const tabList: CardTabListType[] = [
    {
        key: "",
        label: (
            <Space>
                <MenuOutlined />
                <Trans>Basic info</Trans>
            </Space>
        ),
    },
    {
        key: "build",
        label: (
            <Space>
                <SettingOutlined />
                <Trans>Build settings</Trans>
            </Space>
        ),
    },
    // {
    //     key: "buildTasks",
    //     label: (
    //         <Space>
    //             <BuildOutlined />
    //             <Trans>Build tasks</Trans>
    //         </Space>
    //     ),
    // },
];

export default function ApplicationDetail() {
    const { id } = useParams();
    const navigate = useNavigate();

    const { pathname } = useLocation();
    const defaultActiveTab = useMemo(() => {
        return getActivePath(`/app/${id}`, pathname);
    }, [pathname, id]);
    console.log("defaultActiveTab", defaultActiveTab)

    const onTabChange = (key: string) => {
        navigate(`/app/${id}/${key}`);
    };

    return (
        <Card
            title={<div className={"text-xl mb-2"}>{id}</div>}
            defaultActiveTabKey={defaultActiveTab}
            tabList={tabList}
            onTabChange={onTabChange}
            tabProps={{ size: "middle" }}
            extra={<DetailExtra />}
        >
            <Outlet />
        </Card>
    );
}