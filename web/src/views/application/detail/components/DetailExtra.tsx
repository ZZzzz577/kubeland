import { Space } from "antd";
import DeleteButton from "@/views/application/common/components/DeleteButton.tsx";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import EditorButton from "@/views/application/common/components/EditorButton.tsx";

export default function DetailExtra() {
    return (
        <Space>
            <EditorButton type={"primary"} icon={<EditOutlined />} />
            <DeleteButton type={"primary"} icon={<DeleteOutlined />} />
        </Space>
    );
}