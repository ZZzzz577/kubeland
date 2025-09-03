import { Button } from "antd";
import { useLocation, useNavigate, useParams } from "react-router";
import { useLingui } from "@lingui/react/macro";
import type { BaseButtonProps } from "antd/es/button/button";
import getActivePath from "@/views/application/detail/utils/tab.ts";
import { Fragment } from "react";

export default function EditButton(props: BaseButtonProps) {
    const { name } = useParams();
    const { t } = useLingui();
    const navigate = useNavigate();
    const { pathname } = useLocation();
    const activePath = getActivePath(`/app/${name}`, pathname);
    const next = `/app/${name}/edit/${activePath}`;

    if (!name) {
        return <Fragment />;
    }
    return (
        <Button {...props} onClick={() => navigate(next)}>
            {t`Edit`}
        </Button>
    );
}