import { Outlet, useLocation, useNavigate, useParams } from "react-router";
import { Trans } from "@lingui/react/macro";
import { Card, Space } from "antd";
import type { CardTabListType } from "antd/es/card/Card";
import { MenuOutlined, SettingOutlined } from "@ant-design/icons";
import getActivePath from "@/views/application/detail/utils/tab.ts";

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
];

export default function ApplicationModify() {
    const { name } = useParams();
    const navigate = useNavigate();

    const location = useLocation();
    const defaultActiveTab = getActivePath(`/app/${name}/edit`, location.pathname);

    const onTabChange = (key: string) => {
        navigate(`/app/${name}/edit/${key}`);
    };

    return (
        <Card
            title={<div className={"text-xl mb-2"}>{name}</div>}
            defaultActiveTabKey={defaultActiveTab}
            tabList={tabList}
            onTabChange={onTabChange}
            tabProps={{ size: "middle" }}
        >
            <Outlet />
        </Card>
    );
}
