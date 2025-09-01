import { Button } from "antd";
import { useLocation, useNavigate, useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import type { BaseButtonProps } from "antd/es/button/button";
import getActivePath from "@/views/application/common/utils/tab.ts";

export default function EditorButton(props: BaseButtonProps) {
    const { id } = useParams();
    const { t } = useLingui();
    const navigate = useNavigate();
    const { pathname } = useLocation();
    const activePath = getActivePath(`/app/${id}`, pathname);
    return (
        <Button {...props} onClick={() => navigate(`/app/${id}/edit/${activePath}`)}>
            {t`Edit`}
        </Button>
    );
}