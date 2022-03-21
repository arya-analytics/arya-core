import { useThemeContext } from "./Theme/Theme";
import { Box, Grid } from "@mui/material";
import { Nav } from "./Nav/Nav";
import { Pages } from "./Page/Pages";

function App() {
  return (
    <Box sx={{ flexGrow: 1, height: "100vh" }}>
      <Grid container sx={{ height: "100%" }}>
        <Nav />
        <Pages />
      </Grid>
    </Box>
  );
}

export default App;
