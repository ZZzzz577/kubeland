import { type ApiV1ClusterConnection, ApiV1ClusterConnectionTypeEnum } from "@/generated";
import { useLingui } from "@lingui/react/macro";
import { Descriptions, type DescriptionsProps, Space } from "antd";
import { useMemo } from "react";
import { getAuthType } from "@/views/cluster/commons/components/AuthType.tsx";
import TestConnection from "@/views/cluster/detail/components/TestConnection.tsx";

export default function ConnectionDescription(props: { connection?: ApiV1ClusterConnection }) {
    const { connection } = props;
    const { t } = useLingui();

    const items = useMemo(() => {
        const items: DescriptionsProps["items"] = [
            {
                label: t`Address`,
                children: connection?.address,
            },
            {
                label: t`CA`,
                children: <div className={"whitespace-pre-wrap"}>{connection?.ca}</div>,
            },
            {
                label: t`Auth type`,
                children: getAuthType(connection?.type),
            },
        ];
        if (connection?.type === ApiV1ClusterConnectionTypeEnum.TlsCert) {
            items.push({
                label: t`Cert`,
                children: <div className={"whitespace-pre-wrap"}>{connection?.cert}</div>,
            });
            items.push({
                label: t`Key`,
                children: <div className={"whitespace-pre-wrap"}>{connection?.key}</div>,
            });
        } else if (connection?.type === ApiV1ClusterConnectionTypeEnum.TlsToken) {
            items.push({
                label: t`Token`,
                children: <div>{connection?.token}</div>,
            });
        }
        return items;
    }, [
        connection?.address,
        connection?.ca,
        connection?.cert,
        connection?.key,
        connection?.token,
        connection?.type,
        t,
    ]);
    return (
        <Space direction={"vertical"} size={"large"}>
            <Descriptions
                title={t`Connection`}
                column={1}
                bordered
                items={items}
                styles={{
                    label: { width: 120, padding: 10 },
                    root: { border: "none", borderRadius: 0 },
                }}
            />
            <TestConnection connection={connection} />
        </Space>
    );
}
