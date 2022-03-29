import { Box, Button, Stack, TextField } from "@mui/material";
import { AryaIcon } from "../../Icons/Arya";
import { ToggleThemeSwitch } from "../../Theme/Theme";
import { BarProgress } from "../../Node/NodeDetail/NodeDetail";

export interface LoginProps {}

export const Login = ({}: LoginProps) => {
  return (
    <Box
      sx={{
        width: "100vw",
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <LoginThemeSwitch />
      <LoginDiagnostics />
      <LoginForm />
    </Box>
  );
};

const LoginForm = () => {
  return (
    <Stack
      direction="column"
      spacing={6}
      sx={{
        display: "flex",
        alignItems: "center",
        marginTop: -7,
      }}
    >
      <AryaIcon fontSize="large" sx={{ fontSize: 150 }} />
      <Stack direction="column" spacing={4}>
        <TextField
          label="Username"
          sx={{ width: 460 }}
          size={"small"}
          variant={"standard"}
        />
        <TextField
          label="Password"
          sx={{ width: 460 }}
          size={"small"}
          variant={"standard"}
        />
      </Stack>
      <Button variant="contained" size="medium">
        Login
      </Button>
    </Stack>
  );
};

const LoginThemeSwitch = () => {
  return (
    <ToggleThemeSwitch
      size="small"
      sx={{
        position: "absolute",
        bottom: "0",
        left: "0",
        zIndex: "1",
        m: "13px",
      }}
    />
  );
};

const LoginDiagnostics = () => {
  return (
    <Stack
      direction="row"
      sx={{ position: "absolute", bottom: 0, right: 0, zIndex: 1, m: "13px" }}
    >
      <BarProgress name="Live Nodes" progress={88} fill="green" width={400} />
    </Stack>
  );
};
