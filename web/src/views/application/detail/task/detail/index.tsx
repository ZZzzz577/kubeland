import LogView from "@/components/logview";
import { useParams } from "react-router";
import { Card } from "antd";

export default function BuildTaskDetail() {
    const { name, task } = useParams();
    const url = `ws://localhost:8000/v1/app/${name}/build/task/${task}/log`;
    return (
        <Card>
            <LogView className={"h-full w-full whitespace-normal"} url={url} />
        </Card>
    );
}
