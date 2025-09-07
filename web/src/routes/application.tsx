import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { AppstoreOutlined } from "@ant-design/icons";
import ApplicationList from "@/views/application/list";
import ApplicationModify from "@/views/application/modify";
import ApplicationDetail from "@/views/application/detail";
import ApplicationCreate from "@/views/application/create";
import BasicInfo from "@/views/application/detail/basic";
import BuildSettings from "@/views/application/detail/build";
import BasicInfoEdit from "@/views/application/modify/basic";
import BuildSettingsEdit from "@/views/application/modify/build";
import BuildTask from "@/views/application/detail/task";
import BuildTaskDetail from "@/views/application/detail/task/detail";

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
                element: <ApplicationCreate />,
            },

            {
                path: ":name",
                name: <Trans>Detail</Trans>,
                element: <ApplicationDetail />,
                children: [
                    {
                        path: "",
                        name: <Trans>Basic</Trans>,
                        element: <BasicInfo />,
                    },
                    {
                        path: "build",
                        name: <Trans>Build</Trans>,
                        element: <BuildSettings />,
                    },
                    {
                        path: "task",
                        name: <Trans>Task</Trans>,
                        element: <BuildTask />,
                    },
                ],
            },

            {
                path: ":name/task/:task",
                element: <BuildTaskDetail />,
            },

            {
                path: ":name/edit",
                name: <Trans>Edit</Trans>,
                element: <ApplicationModify />,
                children: [
                    {
                        path: "",
                        name: <Trans>Basic</Trans>,
                        element: <BasicInfoEdit />,
                    },
                    {
                        path: "build",
                        name: <Trans>Build</Trans>,
                        element: <BuildSettingsEdit />,
                    },
                ],
            },
        ],
    };
};