import { Box, Grid } from "@mui/material";
import { Login } from "./Pages/Login/Login";
import { Nav } from "./Nav/Nav";
import { Pages } from "./Pages/Pages";

function App() {
  return (
    <Box sx={{ flexGrow: 1, height: "100vh" }}>
      <Login />
      {/*<Grid container sx={{ height: "100%" }}>*/}
      {/*  <Nav />*/}
      {/*  <Pages />*/}
      {/*</Grid>*/}
    </Box>
  );
}

export default App;
