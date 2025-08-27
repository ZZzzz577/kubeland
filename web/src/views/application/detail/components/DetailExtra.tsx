import { Space } from "antd";
import DeleteButton from "@/views/application/common/components/DeleteButton.tsx";
import { useParams } from "react-router";
import { DeleteOutlined } from "@ant-design/icons";

export default function DetailExtra() {
    const { id } = useParams();
    return (
        <Space>
            <DeleteButton id={id} type={"primary"} icon={<DeleteOutlined />} />
        </Space>
    );
}