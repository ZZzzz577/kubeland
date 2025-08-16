import { ApiV1ClusterConnectionTypeEnum } from "@/generated";
import { Tag } from "antd";
import { Trans } from "@lingui/react/macro";

export function getAuthType(type?: ApiV1ClusterConnectionTypeEnum) {
    switch (type) {
        case ApiV1ClusterConnectionTypeEnum.TlsCert:
            return (
                <Tag color={"blue"}>
                    <Trans>CERT</Trans>
                </Tag>
            );
        case ApiV1ClusterConnectionTypeEnum.TlsToken:
            return (
                <Tag color={"blue"}>
                    <Trans>CERT</Trans>
                </Tag>
            );
        default:
            return "";
    }
}
