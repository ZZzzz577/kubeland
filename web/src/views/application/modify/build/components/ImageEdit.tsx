import { Form, Select, Tag } from "antd";
import { Title } from "@/views/application/modify/build";
import { useLingui } from "@lingui/react/macro";
import { useRequest } from "ahooks";
import { imageApi } from "@/api";

export default function ImageEdit() {
    const {t} = useLingui();
    const {Item} = Form;
    const {loading, data} = useRequest(imageApi.imageServiceListImageRepos.bind(imageApi));
    return (
        <>
            <Title>{t`Image repository settings`}</Title>

            <Item name={["image", "repoName"]} label={t`Repo URL`} rules={[{required: true}]}>
                <Select
                    className={"min-w-50"}
                    loading={loading}
                    options={data?.items?.map((item) => ({
                        value: item.name,
                        label: <><Tag color={"blue"}>{item.name}</Tag>{item.url}</>
                    }))}/>
            </Item>
        </>
    );
}