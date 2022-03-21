import { GrainOutlined, NoiseAwareOutlined } from "@mui/icons-material";
import { Box, Grid, Stack, Tabs, Tab, Typography } from "@mui/material";
import { motion } from "framer-motion";
import React, { useRef } from "react";
import {
  Hexagon,
  HexagonProgress,
  MultiHexagonProgress,
} from "../NodeMap/Hexagon";
import {NodeMap, NodeMapPage} from "../NodeMap/NodeMap";

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

export const Page = () => {
  const constraintsRef = useRef(null);
  const [value, setValue] = React.useState(0);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  return (
    <Grid item container xs direction="column">
      <Grid
        height={42}
        borderBottom="1px solid"
        borderColor="divider"
        paddingLeft={1}
        display="flex"
        alignItems="center"
        justifyContent="space-between"
      >
        <Stack direction="row" alignItems="center" spacing={1}>
          <GrainOutlined sx={{ color: "text.primary" }} />
          <Typography variant="h6">Cluster</Typography>
        </Stack>
      </Grid>
      <Grid xs container item>
        <NodeMapPage />
        <Grid xs>
          <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
            <Tabs
              value={value}
              onChange={handleChange}
              aria-label="basic tabs example"
            >
              <Tab label="Node Details" />
              <Tab label="Storage" />
              <Tab label="Alerts" />
            </Tabs>
          </Box>
          <TabPanel value={value} index={0}>
            Item One
          </TabPanel>
          <TabPanel value={value} index={1}>
            Item Two
          </TabPanel>
          <TabPanel value={value} index={2}>
            Item Three
          </TabPanel>
        </Grid>
      </Grid>
    </Grid>
  );
};
