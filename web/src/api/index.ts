import { ClusterServiceApi, Configuration } from "@/generated";

const config = new Configuration({
  basePath: "http://localhost:8000",
});

export const clusterApi = new ClusterServiceApi(config)