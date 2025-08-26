import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { AppstoreOutlined } from "@ant-design/icons";
import ApplicationList from "@/views/application/list";
import ApplicationModify from "@/views/application/modify";

export const Application = (): Route => {
    return {
        path: "/app/application",
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
                name: <Trans>create</Trans>,
                element: <ApplicationModify />,
            },
            {
                path: ":id/edit",
                name: <Trans>create</Trans>,
                element: <ApplicationModify />,
            },
        ],
    };
};