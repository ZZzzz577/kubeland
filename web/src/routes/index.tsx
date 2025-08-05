import type { ReactNode } from "react";
import type { RouteObject } from "react-router";
import { Home } from "@/routes/home";

interface RouteMenuConfig {
  icon?: ReactNode;
}

export type Route = RouteObject & {
  name?: ReactNode;
  menu?: RouteMenuConfig;
  children?: Route[];
};

const routers: Route[] = [Home()];
export default routers;
