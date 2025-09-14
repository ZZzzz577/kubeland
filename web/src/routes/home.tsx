import type { Route } from "@/routes/index";
import { HomeOutlined } from "@ant-design/icons";
import AppLayout from "@/layout";
import { Trans } from "@lingui/react/macro";

export const Home: Route = {
    path: "/",
    element: <AppLayout />,
    name: <Trans>home</Trans>,
    menu: {
        icon: <HomeOutlined />,
    },
    children: [
        {
            path: "",
            element: <div>home</div>,
        },
    ],
};
