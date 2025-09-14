import {
    ApplicationServiceApi,
    BuildSettingsServiceApi,
    BuildTaskServiceApi,
    ClusterServiceApi,
    Configuration,
    GitServiceApi,
    ImageServiceApi,
} from "@/generated";
import { StatusMiddleware } from "@/api/middleware/status.ts";

const config = new Configuration({
    basePath: "http://localhost:8000",
    middleware: [StatusMiddleware],
});

export const clusterApi = new ClusterServiceApi(config);
export const applicationApi = new ApplicationServiceApi(config);
export const buildSettingsApi = new BuildSettingsServiceApi(config);
export const buildTaskApi = new BuildTaskServiceApi(config);
export const gitApi = new GitServiceApi(config);
export const imageApi = new ImageServiceApi(config);
