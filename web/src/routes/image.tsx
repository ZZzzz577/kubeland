import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { DockerOutlined } from "@ant-design/icons";
import ImageRepoList from "@/views/image/list";
import ImageRepoCreate from "@/views/image/create";
import ImageRepoDetail from "@/views/image/detail";
import ImageRepoUpdate from "@/views/image/update";

export const Image: Route = {
    path: "/image",
    element: <AppLayout />,
    name: <Trans>Image Repository</Trans>,
    menu: {
        icon: <DockerOutlined />,
    },
    children: [
        {
            path: "",
            element: <ImageRepoList />,
        },
        {
            path: "create",
            name: <Trans>Create</Trans>,
            element: <ImageRepoCreate />,
        },
        {
            path: ":name",
            name: <Trans>Detail</Trans>,
            element: <ImageRepoDetail />,
        },
        {
            path: ":name/edit",
            name: <Trans>Edit</Trans>,
            element: <ImageRepoUpdate />,
        },
    ],
};