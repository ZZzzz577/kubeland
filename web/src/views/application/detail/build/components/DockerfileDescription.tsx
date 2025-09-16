import { Descriptions } from "antd";
import type { DescriptionsItemType } from "antd/es/descriptions";

export default function DockerfileDescription(props: {
    dockerfile?: string
}) {
    const { dockerfile } = props;

    const items: DescriptionsItemType[] = [
        {
            label: "",
            children: <div className={"w-full min-h-96 border p-2 border-gray-300"}>{dockerfile}</div>
        }
    ];
    return (
        <Descriptions
            title={"Dockerfile"}
            items={items}
        />

    );
}