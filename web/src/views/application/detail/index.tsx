import { Card, Menu, type MenuProps, Space } from "antd";
import { Trans } from "@lingui/react/macro";
import { useMemo } from "react";
import {
    BuildOutlined,
    DeleteOutlined,
    DockerOutlined,
    GitlabOutlined,
    MenuOutlined,
    SettingOutlined,
} from "@ant-design/icons";
import { Outlet, useLocation, useNavigate, useParams } from "react-router";
import getActivePath from "@/views/application/detail/utils/tab.ts";
import DeleteButton from "@/views/application/detail/components/DeleteButton.tsx";

type MenuItem = Required<MenuProps>["items"][number];
const menuItems: MenuItem[] = [
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
    {
        key: "git",
        label: (
            <Space>
                <GitlabOutlined />
                <Trans>Git repo</Trans>
            </Space>
        ),
    },
    {
        key: "task",
        label: (
            <Space>
                <BuildOutlined />
                <Trans>Build tasks</Trans>
            </Space>
        ),
    },
    {
        key: "image",
        label: (
            <Space>
                <DockerOutlined />
                <Trans>Image repo</Trans>
            </Space>
        ),
    },
];

export default function ApplicationDetail() {
    const { name } = useParams();
    const navigate = useNavigate();

    const { pathname } = useLocation();
    const defaultActiveTab = useMemo(() => {
        return getActivePath(`/app/${name}`, pathname);
    }, [pathname, name]);

    const onClick: MenuProps["onClick"] = (e) => {
        navigate(`/app/${name}/${e.key}`);
    };

    return (
        <>
            <Card
                styles={{ body: { padding: 0 } }}
                title={<div className={"text-3xl"}>{name}</div>}
                extra={
                    <Space>
                        <DeleteButton icon={<DeleteOutlined />} size={"large"} danger name={name} />
                    </Space>
                }
            >
                <Menu mode="horizontal" selectedKeys={[defaultActiveTab]} onClick={onClick} items={menuItems} />
            </Card>
            <Outlet />
        </>
    );
}