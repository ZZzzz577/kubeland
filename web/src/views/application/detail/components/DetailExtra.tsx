import { Space } from "antd";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import EditButton from "@/views/application/detail/components/EditButton.tsx";
import DeleteButton from "@/views/application/detail/components/DeleteButton.tsx";

export default function DetailExtra() {
    return (
        <Space>
            <EditButton type={"primary"} icon={<EditOutlined />} />
            <DeleteButton type={"primary"} icon={<DeleteOutlined />} />
        </Space>
    );
}