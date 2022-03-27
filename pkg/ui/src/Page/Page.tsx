import { Grid, Stack, SxProps, Typography } from "@mui/material";
import { cloneElement, PropsWithChildren, ReactElement } from "react";

interface PageHeadingProps {
  icon: ReactElement;
  variant: "heading" | "subheading";
  children: string;
}

const ICON_SIZE_HEADING_LARGE = 24;
const ICON_SIZE_HEADING_SMALL = 18;

export const PageHeading = ({ icon, variant, children }: PageHeadingProps) => {
  const iconProps = {
    sx: {
      color: "text.primary",
      fontSize:
        variant == "heading"
          ? ICON_SIZE_HEADING_LARGE
          : ICON_SIZE_HEADING_SMALL,
    },
  };
  return (
    <Grid
      item
      height={36}
      borderBottom="1px solid"
      borderColor="divider"
      paddingLeft={1}
      display="flex"
      alignItems="center"
      justifyContent="space-between"
    >
      <Stack direction="row" alignItems="center" spacing={1}>
        {cloneElement(icon, iconProps)}
        <Typography variant={variant == "heading" ? "h6" : "subtitle1"}>
          {children}
        </Typography>
      </Stack>
    </Grid>
  );
};

export interface PageProps extends PropsWithChildren<any> {
  xs?: number | "auto" | boolean;
  direction?: "row" | "column";
  sx: SxProps;
}

export const Page = ({
  children,
  xs = true,
  direction = "column",
  ...props
}: PageProps) => {
  return (
    <Grid item container xs={xs} direction={direction} {...props}>
      {children}
    </Grid>
  );
};
