import { ClusterServiceApi, Configuration } from "@/generated";
import { StatusMiddleware } from "@/api/middleware/status.ts";

const config = new Configuration({
    basePath: "http://localhost:8000",
    middleware: [StatusMiddleware],
});

export const clusterApi = new ClusterServiceApi(config);
