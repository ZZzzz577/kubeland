import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { AppstoreOutlined } from "@ant-design/icons";
import ApplicationList from "@/views/application/list";
import ApplicationModify from "@/views/application/modify";
import ApplicationDetail from "@/views/application/detail";

export const Application = (): Route => {
    return {
        path: "/app",
        element: <AppLayout />,
        name: <Trans>Application</Trans>,
        menu: {
            icon: <AppstoreOutlined />,
        },
        children: [
            {
                path: "",
                name: <Trans>Application</Trans>,
                menu: {},
                element: <ApplicationList />,
            },
            {
                path: "create",
                name: <Trans>Create</Trans>,
                element: <ApplicationModify />,
            },
            {
                path: ":id/edit",
                name: <Trans>Edit</Trans>,
                element: <ApplicationModify />,
            },
            {
                path: ":id/detail",
                name: <Trans>Detail</Trans>,
                element: <ApplicationDetail />,
            },
        ],
    };
};