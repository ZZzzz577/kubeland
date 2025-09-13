import { Button } from "antd";
import { useNavigate } from "react-router";
import { useLingui } from "@lingui/react/macro";
import type { BaseButtonProps } from "antd/es/button/button";
import { Fragment } from "react";

export default function EditButton(props: { name?: string } & BaseButtonProps) {
    const { name } = props;
    const { t } = useLingui();
    const navigate = useNavigate();
    const next = `/app/${name}/edit`;

    if (!name) {
        return <Fragment />;
    }
    return (
        <Button {...props} onClick={() => navigate(next)}>
            {t`Edit`}
        </Button>
    );
}
