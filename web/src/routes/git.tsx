import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { GitlabOutlined } from "@ant-design/icons";
import GitRepoList from "@/views/git/list";
import GitRepoCreate from "@/views/git/create";
import GitRepoUpdate from "@/views/git/update";
import GitRepoDetail from "@/views/git/detail";

export const Git: Route = {
    path: "/git",
    element: <AppLayout />,
    name: <Trans>Git Repository</Trans>,
    menu: {
        icon: <GitlabOutlined />,
    },
    children: [
        {
            path: "",
            element: <GitRepoList />,
        },
        {
            path: "create",
            name: <Trans>Create</Trans>,
            element: <GitRepoCreate />,
        },
        {
            path: ":name",
            name: <Trans>Detail</Trans>,
            element: <GitRepoDetail />,
        },
        {
            path: ":name/edit",
            name: <Trans>Edit</Trans>,
            element: <GitRepoUpdate />,
        },
    ],
};