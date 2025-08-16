import type { Route } from "@/routes/index.tsx";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";
import { KubernetesOutlined } from "@ant-design/icons";
import ClusterModify from "@/views/cluster/modify";
import ClusterList from "@/views/cluster/list";
import ClusterDetail from "@/views/cluster/detail";

export const Cluster = (): Route => {
    return {
        path: "/cluster",
        element: <AppLayout />,
        name: <Trans>cluster</Trans>,
        menu: {
            icon: <KubernetesOutlined />,
        },
        children: [
            {
                path: "",
                element: <ClusterList />,
            },
            {
                path: ":id/detail",
                element: <ClusterDetail />,
            },
            {
                path: "create",
                name: <Trans>create</Trans>,
                element: <ClusterModify />,
            },
            {
                path: ":id/edit",
                element: <ClusterModify />,
            },
        ],
    };
};