import { useFileUpload } from "@/components/file/hooks.ts";
import { Trans, useLingui } from "@lingui/react/macro";
import { Alert, Button, Flex, Form, Select, Space, Tag, Upload } from "antd";
import { clusterApi } from "@/api";
import { useRequest } from "ahooks";
import { type ApiV1ClusterCluster, ApiV1ClusterConnectionTypeEnum } from "@/generated";
import { UploadOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import useApp from "antd/es/app/useApp";

export default function UploadKubeConfig() {
    const { t } = useLingui();
    const { notification } = useApp();
    const form = Form.useFormInstance<ApiV1ClusterCluster>();
    const { Item } = Form;
    const [currentCtxName, setCurrentCtxName] = useState<string>();

    const { data, run } = useRequest(clusterApi.clusterServiceResolveKubeConfig.bind(clusterApi), {
        manual: true,
        onSuccess: (data) => {
            if (!data.items?.length) {
                notification.error({
                    message: t`resolve kubeconfig error`,
                    description: t`No valid context. Please check your kubeconfig file`,
                });
                return;
            }
            const currentCtx = data.items?.find((item) => item.current);
            if (currentCtx?.name) {
                setCurrentCtxName(currentCtx.name);
            }
        },
        onError: (error) => {
            notification.error({
                message: t`resolve kubeconfig error`,
                description: error.message,
            });
        },
    });
    useEffect(() => {
        if (data?.items?.length && currentCtxName) {
            const currentCtx = data.items.find((item) => item.name === currentCtxName);
            if (currentCtx) {
                form.setFieldsValue({
                    connection: {
                        address: currentCtx.cluster?.server,
                        type: ApiV1ClusterConnectionTypeEnum.TlsCert,
                        ca: currentCtx.cluster?.ca,
                        cert: currentCtx.user?.cert,
                        key: currentCtx.user?.key,
                    },
                });
            }
        }
    }, [data, currentCtxName, form]);

    const uploadKubeConfig = useFileUpload({
        onSuccess: (_, content) => {
            run({
                apiV1ClusterResolveKubeConfigRequest: {
                    content: content,
                },
            });
        },
        onError: (error) => {
            notification.error({
                message: t`failed to upload kubeconfig file`,
                description: error.message,
            });
        },
    });
    return (
        <>
            <Flex className={"!mb-5"} align={"center"}>
                <div className={"w-25 text-base font-medium"}>{t`Connection`}</div>
                <Alert
                    type={"info"}
                    message={
                        <Trans>
                            <Upload
                                beforeUpload={uploadKubeConfig.handleFileUpload}
                                showUploadList={false}
                                disabled={uploadKubeConfig.loading}
                            >
                                <Button
                                    type={"link"}
                                    size={"small"}
                                    loading={uploadKubeConfig.loading}
                                    icon={<UploadOutlined />}
                                >{t`upload kubeconfig file`}</Button>
                            </Upload>
                            to fill connection settings automatically
                        </Trans>
                    }
                />
            </Flex>
            {(data?.items?.length ?? 0) > 1 && (
                <Item label={"context"}>
                    <Select
                        className={"min-w-100"}
                        value={currentCtxName}
                        onChange={(value) => {
                            setCurrentCtxName(value);
                        }}
                        options={data?.items?.map((item) => ({
                            label: (
                                <Space>
                                    <Tag color={"blue"}>{item.name}</Tag>
                                    <div> {item.cluster?.server}</div>
                                </Space>
                            ),
                            value: item.name,
                        }))}
                    />
                </Item>
            )}
        </>
    );
}
