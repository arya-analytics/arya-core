// @flow
import { ApiOutlined, DashboardOutlined } from "@mui/icons-material";
import {
  Box,
  Divider,
  Grid,
  List,
  ListItem,
  ListItemIcon,
  Stack,
  Tooltip,
} from "@mui/material";
import { AryaIcon } from "../Icons/Arya";
import { ToggleThemeSwitch } from "../Theme/Theme";

interface Props {}

export const Nav = (props: Props) => {
  return (
    <Grid
      item
      width={60}
      borderRight="1px solid"
      borderColor="divider"
      paddingTop={2}
      paddingBottom={2}
      display="flex"
      justifyContent="space-between"
      direction="column"
      alignItems="center"
    >
      <NavTop />
      <NavBottom />
    </Grid>
  );
};

const NavTop = () => {
  return (
    <Stack spacing={2} width="100%">
      <NavHeader />
      <Divider />
      <NavMenu />
    </Stack>
  );
};

const NavBottom = () => {
  return <ToggleThemeSwitch size="small" />;
};

const NavHeader = () => {
  return (
    <Box sx={{ textAlign: "center" }}>
      <AryaIcon fontSize="large" />
    </Box>
  );
};

const navMenuButtons = [
  {
    name: "Dashboard",
    icon: <DashboardOutlined />,
  },
  {
    name: "Devices",
    icon: <ApiOutlined />,
  },
];

const NavMenu = () => {
  return (
    <Box role="presentation">
      <List>
        {navMenuButtons.map(({ name, icon }) => {
          return (
            <Tooltip title={name} placement="right">
              <ListItem button key={name}>
                <ListItemIcon>{icon}</ListItemIcon>
              </ListItem>
            </Tooltip>
          );
        })}
      </List>
    </Box>
  );
};
