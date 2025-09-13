import { Card, Space } from "antd";
import DockerfileDescription from "@/views/application/detail/build/components/DockerfileDescription.tsx";
import GitDescription from "@/views/application/detail/build/components/GitDescription.tsx";
import ImageDescription from "@/views/application/detail/build/components/ImageDescription.tsx";


export default function BuildSettings() {

    return (
        <Card
            className={"!rounded-t-none"}
            variant="borderless"
        >
            <Space direction={"vertical"} size={"large"}>
                <GitDescription />
                <ImageDescription />
                <DockerfileDescription />
            </Space>
        </Card>
    );

}