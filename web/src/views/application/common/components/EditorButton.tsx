import { Button } from "antd";
import { useNavigate } from "react-router";
import { useLingui } from "@lingui/react/macro";
import type { BaseButtonProps } from "antd/es/button/button";

export default function EditorButton(props: { id?: string } & BaseButtonProps) {
    const { id } = props;
    const { t } = useLingui();
    const navigate = useNavigate();
    return (
        <Button {...props} onClick={() => navigate(`/app/${id}/edit`)}>
            {t`Edit`}
        </Button>
    );
}