import type { ReactNode } from "react";
import type { RouteObject } from "react-router";
import { Home } from "@/routes/home";
import { Cluster } from "@/routes/cluster.tsx";
import { Application } from "@/routes/application.tsx";
import { Image } from "@/routes/image.tsx";

interface RouteMenuConfig {
    icon?: ReactNode;
}

export type Route = RouteObject & {
    name?: ReactNode;
    menu?: RouteMenuConfig;
    children?: Route[];
};

const routes: Route[] = [Home, Cluster, Image, Application];
export default routes;
