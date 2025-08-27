import { Descriptions, Space } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useParams } from "react-router";
import { ApplicationCluster } from "@/views/application/list/components/ApplicationListTable.tsx";
import { EditOutlined } from "@ant-design/icons";
import EditorButton from "@/views/application/common/components/EditorButton.tsx";
import type { ApiV1ApplicationApplication } from "@/generated";

export default function BasicInfo(props: { info?: ApiV1ApplicationApplication }) {
    const { t } = useLingui();
    const { id } = useParams();
    const { info } = props;
    const { Item } = Descriptions;

    return (
        <Descriptions
            title={<div className={"text-base"}>{t`Basic info`}</div>}
            styles={{ header: { marginBottom: 10 }, label: { width: 150 } }}
            column={2}
            bordered
            extra={
                <Space>
                    <EditorButton id={id} type={"primary"} icon={<EditOutlined />} />
                </Space>
            }
        >
            <Item label={t`Name`} span={"filled"}>
                {info?.name}
            </Item>
            <Item label={t`Cluster`} span={"filled"}>
                <ApplicationCluster id={info?.clusterId} />
            </Item>
            <Item label={t`Description`} span={"filled"}>
                {info?.description}
            </Item>
            <Item label={t`Create time`}>{info?.createdAt?.toLocaleString()}</Item>
            <Item label={t`Update time`}>{info?.updateAt?.toLocaleString()}</Item>
        </Descriptions>
    );
}