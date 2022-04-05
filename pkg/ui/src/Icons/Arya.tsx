import {SvgIcon, SvgIconProps, useTheme} from "@mui/material";

import {ReactComponent as White} from "../media/icon-white.svg"
import {ReactComponent as FullWhite} from "../media/icon-full-title-white.svg"
import {ReactComponent as Gradient} from "../media/icon-gradient.svg"
import {ReactComponent as FullGradient} from "../media/icon-full-title.svg"

import {FunctionComponent, SVGProps} from "react";

interface ThemedIconProps extends SvgIconProps {
    light: FunctionComponent<SVGProps<SVGSVGElement>>,
    dark: FunctionComponent<SVGProps<SVGSVGElement>>,
}

const ThemedIcon = (props: ThemedIconProps) => {
    const {palette: {mode}} = useTheme();
    const icon = mode == "light" ? props.light : props.dark;
    return <SvgIcon component={icon} inheritViewBox {...props}  />
}

export const AryaIcon = (props: SvgIconProps) => {
    return <ThemedIcon light={Gradient} dark={White} {...props} />
}

export const AryaFullIcon = (props: SvgIconProps) => {
    return <ThemedIcon light={FullGradient} dark={FullWhite}/>
}