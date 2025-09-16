import { Button } from "antd";
import { useNavigate, useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import { PlusCircleOutlined } from "@ant-design/icons";

export default function CreateBuildTask() {
    const { name } = useParams();
    const { t } = useLingui();
    const navigate = useNavigate();

    return (
        <Button
            type={"primary"}
            size={"middle"}
            icon={<PlusCircleOutlined />}
            onClick={() => {
                navigate(`/app/${name}/git`);
            }}
        >{t`Create`}</Button>
    );
}