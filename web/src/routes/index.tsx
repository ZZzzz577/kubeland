import type { ReactNode } from "react";
import type { RouteObject } from "react-router";
import { Home } from "@/routes/home";
import { Cluster } from "@/routes/cluster.tsx";

interface RouteMenuConfig {
    icon?: ReactNode;
}

export type Route = RouteObject & {
    name?: ReactNode;
    menu?: RouteMenuConfig;
    children?: Route[];
};

const routers: Route[] = [Home(), Cluster()];
export default routers;
