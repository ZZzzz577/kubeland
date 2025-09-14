import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { DockerOutlined } from "@ant-design/icons";
import ImageRepoList from "@/views/image/list";
import ImageRepositoryCreate from "@/views/image/create";

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
            element: <ImageRepositoryCreate />,
        },
    ],
};
