import { DoubleLeftOutlined } from "@ant-design/icons";
import { Button } from "antd";
import { useLingui } from "@lingui/react/macro";
import { useNavigate } from "react-router";
import { useCallback } from "react";

export default function BackButton(props: { url?: string }) {
    const { url } = props;
    const { t } = useLingui();
    const navigate = useNavigate();
    const back = useCallback(() => {
        if (url) {
            navigate(url);
        } else {
            navigate(-1);
        }
    }, [navigate, url]);
    return <Button icon={<DoubleLeftOutlined />} onClick={back}>{t`Back`}</Button>;
}