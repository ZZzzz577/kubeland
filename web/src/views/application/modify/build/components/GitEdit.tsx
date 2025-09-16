import { Form, Input, Select, Tag } from "antd";
import { Title } from "@/views/application/modify/build";
import { useLingui } from "@lingui/react/macro";
import { useRequest } from "ahooks";
import { gitApi } from "@/api";

export default function GitEdit() {
    const { t } = useLingui();
    const { Item } = Form;
    const { loading, data } = useRequest(gitApi.gitServiceListGitRepos.bind(gitApi));
    return (
        <>
            <Title>{t`Git repository settings`}</Title>

            <Item name={["git", "repoName"]} label={t`Repo URL`} rules={[{ required: true }]}>
                <Select
                    className={"min-w-50"}
                    loading={loading}
                    options={data?.items?.map((item) => ({
                        value: item.name,
                        label: <><Tag color={"blue"}>{item.name}</Tag>{item.url}</>
                    }))} />
            </Item>
            <Item name={["git", "repoPath"]} label={t`Repo Path`} rules={[{ required: true }]}>
                <Input />
            </Item>
        </>
    );
}