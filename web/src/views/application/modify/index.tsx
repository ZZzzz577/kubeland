import { Outlet, useLocation, useNavigate, useParams } from "react-router";
import { Trans } from "@lingui/react/macro";
import { Card, Space } from "antd";
import { useMemo } from "react";
import type { CardTabListType } from "antd/es/card/Card";
import { MenuOutlined, SettingOutlined } from "@ant-design/icons";

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
    const { id } = useParams();
    const navigate = useNavigate();

    const location = useLocation();
    const defaultActiveTab = useMemo(() => {
        const segments = location.pathname.split("/").filter(Boolean);
        return segments.length > 0 ? segments[segments.length - 1] : "";
    }, [location]);

    const onTabChange = (key: string) => {
        navigate(`/app/${id}/edit/${key}`);
    };

    return (
        <Card
            title={<div className={"text-xl mb-2"}>{id}</div>}
            defaultActiveTabKey={defaultActiveTab}
            tabList={tabList}
            onTabChange={onTabChange}
            tabProps={{ size: "middle" }}
        >
            <Outlet />
        </Card>
    );
}